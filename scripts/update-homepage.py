"""Refresh only the README download block from an actually published release."""
import argparse
import datetime
import json
import re
import urllib.request
from pathlib import Path

REPOSITORY = "vertex-ai-llc/easyai-downloads"
BASE = f"https://github.com/{REPOSITORY}"
START, END = "<!-- release-downloads:start -->", "<!-- release-downloads:end -->"

def download_block(release):
    tag = release["tag_name"]
    if not re.fullmatch(r"v[0-9]+.[0-9]+.[0-9]+", tag) or release.get("draft") or release.get("prerelease"):
        raise ValueError("Expected a published stable release")
    version = tag[1:]
    published = datetime.datetime.fromisoformat(release["published_at"].replace("Z", "+00:00"))
    date = published.astimezone(datetime.timezone(datetime.timedelta(hours=8))).date()
    assets = {asset["name"]: asset for asset in release["assets"] if asset.get("size", 0) > 0}
    options = [
        ("Windows-x64-Setup.exe", "windows", "Windows 安装包"),
        ("macOS-arm64.dmg", "mac-arm", "Mac 苹果芯片版"),
        ("macOS-x64.dmg", "mac-intel", "Mac Intel 版"),
    ]
    buttons = []
    for suffix, artwork, label in options:
        name = f"EasyAI-{version}-{suffix}"
        asset = assets.get(name)
        if not asset:
            buttons.append(f'<a href="{BASE}/releases">{label}：查看可用版本</a>')
            continue
        url = f"{BASE}/releases/download/{tag}/{name}"
        if asset["browser_download_url"] != url:
            raise ValueError("Release asset URL does not match this repository")
        buttons.append(f'<a href="{url}"><img src="assets/download-{artwork}.svg" alt="下载 {label}" width="260" height="64"></a>')
    return "\n".join([
        START, f"**当前版本：{tag}** · {date} 发布", "",
        '<p align="center">' + chr(10).join(buttons) + '</p>', "",
        "Windows 10 / 11 · 64 位　｜　macOS · 苹果 M 系列 / Intel", "",
        f"[查看这个版本的更新说明]({BASE}/releases/tag/{tag}) · [全部下载文件]({BASE}/releases/latest)", END,
    ])

def update(readme, release):
    text = readme.read_text(encoding="utf-8")
    if text.count(START) != 1 or text.count(END) != 1:
        raise ValueError("README must contain exactly one download block")
    before, rest = text.split(START)
    _, after = rest.split(END)
    changed = before + download_block(release) + after
    if changed != text:
        readme.write_text(changed, encoding="utf-8")
    return changed != text

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--release-json", type=Path)
    parser.add_argument("--readme", type=Path, default=Path(__file__).resolve().parents[1] / "README.md")
    args = parser.parse_args()
    if args.release_json:
        release = json.loads(args.release_json.read_text(encoding="utf-8"))
    else:
        request = urllib.request.Request(f"https://api.github.com/repos/{REPOSITORY}/releases/latest", headers={"User-Agent": "EasyAI-downloads-homepage", "Accept": "application/vnd.github+json"})
        with urllib.request.urlopen(request, timeout=30) as response:
            release = json.load(response)
    print("Homepage downloads updated" if update(args.readme, release) else "Homepage downloads already current")
