"""Mirror the latest published GitHub release; promote verified files to stable R2 names."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import time
import urllib.request
from datetime import datetime, timezone

REPO = "vertex-ai-llc/easyai-downloads"
PUBLIC = "https://easyai-download.postlyn.com"
TARGETS = {"windows-x64": "Windows-x64-Setup.exe", "macos-arm64": "macOS-arm64.dmg", "macos-x64": "macOS-x64.dmg"}
STABLE_CACHE = "public, max-age=0, must-revalidate"

def version_tuple(tag):
    if not re.fullmatch(r"v[0-9]+[.][0-9]+[.][0-9]+", tag):
        raise ValueError("Only stable version tags can be mirrored")
    return tuple(map(int, tag[1:].split(".")))

def github_latest():
    headers = {"User-Agent": "EasyAI-R2-Publisher", "Accept": "application/vnd.github+json"}
    if os.environ.get("GH_TOKEN"):
        headers["Authorization"] = "Bearer " + os.environ["GH_TOKEN"]
    with urllib.request.urlopen(urllib.request.Request(f"https://api.github.com/repos/{REPO}/releases/latest", headers=headers), timeout=30) as response:
        return json.load(response)

def release_assets(release):
    tag = release["tag_name"]
    version_tuple(tag)
    if release.get("draft") or release.get("prerelease"):
        raise ValueError("Cannot publish a draft or prerelease")
    assets = {item["name"]: item for item in release["assets"]}
    result = {}
    for platform, suffix in TARGETS.items():
        name = f"EasyAI-{tag[1:]}-{suffix}"
        if name not in assets:
            continue
        item = assets[name]
        expected_url = f"https://github.com/{REPO}/releases/download/{tag}/{name}"
        digest = item.get("digest", "")
        if item.get("size", 0) <= 0 or item.get("browser_download_url") != expected_url or not re.fullmatch(r"sha256:[a-f0-9]{64}", digest):
            raise ValueError(f"Invalid published asset metadata: {name}")
        result[platform] = {"name": name, "version": tag, "size": item["size"], "sha256": digest[7:], "githubUrl": expected_url,
                            "key": f"easyai/EasyAI-{suffix}", "archiveKey": f"easyai/releases/{tag}/{name}"}
        result[platform]["url"] = PUBLIC + "/" + result[platform]["key"]
        if isinstance(item.get("id"), int) and item["id"] > 0:
            result[platform]["assetApiUrl"] = f"https://api.github.com/repos/{REPO}/releases/assets/{item['id']}"
    if not result:
        raise ValueError("Release has no supported installers")
    return result

def file_sha(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()

def download(item, directory):
    path = directory / item["name"]
    if path.exists() and path.stat().st_size == item["size"] and file_sha(path) == item["sha256"]:
        return path
    part = path.with_suffix(path.suffix + ".part")
    print("Downloading " + item["name"], flush=True)
    sources = []
    api_url = item.get("assetApiUrl", "")
    if re.fullmatch(r"https://api[.]github[.]com/repos/" + re.escape(REPO) + r"/releases/assets/[0-9]+", api_url):
        headers = ["Accept: application/octet-stream", "User-Agent: EasyAI-R2-Publisher"]
        if os.environ.get("GH_TOKEN"):
            headers.append("Authorization: Bearer " + os.environ["GH_TOKEN"])
        sources.append((api_url, headers))
    sources.append((item["githubUrl"], []))
    for url, headers in sources:
        # Supply credentials through stdin, never process arguments or logs.
        # curl strips Authorization on cross-host redirects (no --location-trusted).
        config = chr(10).join("header = " + json.dumps(header) for header in headers)
        command = ["curl.exe" if os.name == "nt" else "curl", "--config", "-", "--fail", "--location", "--retry", "2", "--retry-all-errors", "--retry-delay", "3", "--connect-timeout", "20", "--max-time", "120", "--silent", "--show-error", "--output", str(part), url + "?download=" + str(time.time_ns())]
        result = subprocess.run(command, input=config, text=True, check=False)
        if result.returncode == 0 and part.exists() and part.stat().st_size == item["size"] and file_sha(part) == item["sha256"]:
            part.replace(path)
            return path
        print("Download failed or hash mismatch; trying remaining official source", flush=True)
    raise ValueError("Could not download a matching GitHub asset: " + item["name"])

def verify_object(client, bucket, key, item):
    response = client.get_object(Bucket=bucket, Key=key)
    stream = response["Body"]
    digest, size = hashlib.sha256(), 0
    try:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk); size += len(chunk)
    finally:
        stream.close()
    if size != item["size"] or digest.hexdigest() != item["sha256"]:
        raise ValueError("R2 readback mismatch: " + key)

def publish(client, bucket, release, items, directory):
    tag = release["tag_name"]
    try:
        response = client.get_object(Bucket=bucket, Key="easyai/latest.json")
        try: previous = json.loads(response["Body"].read())
        finally: response["Body"].close()
    except client.exceptions.NoSuchKey:
        previous = {}
    if previous and version_tuple(previous["version"]) > version_tuple(tag):
        raise ValueError("Refusing to replace a newer R2 release")
    # Stage and read back every available platform before changing any stable URL.
    for platform, item in items.items():
        metadata = {"version": tag, "sha256": item["sha256"], "github-release": str(release["id"])}
        mime = "application/x-apple-diskimage" if item["name"].endswith(".dmg") else "application/octet-stream"
        item["contentType"] = mime
        client.upload_file(str(directory / item["name"]), bucket, item["archiveKey"], ExtraArgs={"ContentType": mime, "CacheControl": "public, max-age=31536000, immutable", "Metadata": metadata})
        verify_object(client, bucket, item["archiveKey"], item)
        print("Staged and verified " + platform, flush=True)
    for platform, item in items.items():
        client.copy_object(Bucket=bucket, Key=item["key"], CopySource={"Bucket": bucket, "Key": item["archiveKey"]}, MetadataDirective="REPLACE",
                           Metadata={"version": tag, "sha256": item["sha256"], "github-release": str(release["id"])}, ContentType=item["contentType"], CacheControl=STABLE_CACHE,
                           ContentDisposition='attachment; filename="' + item["key"].split("/")[-1] + '"')
        verify_object(client, bucket, item["key"], item)
        print("Published and verified " + item["url"], flush=True)
    # Retain the last working download for platforms absent from a partial release.
    platforms = {**previous.get("platforms", {}), **items}
    manifest = {"version": tag, "releaseId": release["id"], "githubUrl": release["html_url"], "publishedAt": datetime.now(timezone.utc).isoformat(), "platforms": platforms}
    client.put_object(Bucket=bucket, Key="easyai/latest.json", Body=(json.dumps(manifest, indent=2) + chr(10)).encode(), ContentType="application/json", CacheControl=STABLE_CACHE)
    return manifest

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--config", type=Path, help="Local private R2 configuration; CI uses R2_UPLOAD_CONFIG")
    parser.add_argument("--directory", type=Path, default=Path("dist/r2-mirror"))
    args = parser.parse_args()
    settings = json.loads(args.config.read_text(encoding="utf-8-sig") if args.config else os.environ["R2_UPLOAD_CONFIG"])
    import boto3
    from botocore.config import Config
    client = boto3.client("s3", endpoint_url=settings["endpoint"], region_name="auto", aws_access_key_id=settings["access_key_id"], aws_secret_access_key=settings["secret_access_key"],
                          config=Config(signature_version="s3v4", connect_timeout=20, read_timeout=90, retries={"max_attempts": 5, "mode": "standard"}, request_checksum_calculation="when_required", response_checksum_validation="when_required"))
    release = github_latest()
    items = release_assets(release)
    args.directory.mkdir(parents=True, exist_ok=True)
    for item in items.values(): download(item, args.directory)
    # A different release may have become latest while the installers were downloading.
    if github_latest()["id"] != release["id"]: raise RuntimeError("Latest release changed; rerun to mirror the new release")
    manifest = publish(client, settings["bucket"], release, items, args.directory)
    (args.directory / "mirror.json").write_text(json.dumps(manifest, indent=2), encoding="utf-8")
    print(json.dumps({"mirroredVersion": release["tag_name"], "platforms": list(items)}), flush=True)

if __name__ == "__main__": main()
