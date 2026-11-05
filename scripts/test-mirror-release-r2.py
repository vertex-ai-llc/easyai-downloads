import copy
import hashlib
import importlib.util
import io
import json
from pathlib import Path
import tempfile
import unittest

spec = importlib.util.spec_from_file_location("mirror", Path(__file__).with_name("mirror-release-r2.py"))
mirror = importlib.util.module_from_spec(spec); spec.loader.exec_module(mirror)

class Missing(Exception): pass
class FakeR2:
    class exceptions: NoSuchKey = Missing
    def __init__(self): self.objects = {}; self.copies = []; self.corrupt = False
    def get_object(self, Bucket, Key):
        if Key not in self.objects: raise Missing()
        return {"Body": io.BytesIO(self.objects[Key])}
    def upload_file(self, path, bucket, key, ExtraArgs):
        self.objects[key] = b"wrong" if self.corrupt else Path(path).read_bytes()
    def copy_object(self, **args):
        self.copies.append(args)
        self.objects[args["Key"]] = self.objects[args["CopySource"]["Key"]]
    def put_object(self, Bucket, Key, Body, **kwargs): self.objects[Key] = Body

class MirrorTests(unittest.TestCase):
    def release(self):
        name = "EasyAI-1.2.3-Windows-x64-Setup.exe"
        return {"id": 123, "tag_name": "v1.2.3", "html_url": "https://github.com/vertex-ai-llc/easyai-downloads/releases/tag/v1.2.3", "assets": [{"name": name, "size": 4, "digest": "sha256:" + hashlib.sha256(b"test").hexdigest(), "browser_download_url": "https://github.com/vertex-ai-llc/easyai-downloads/releases/download/v1.2.3/" + name}]}
    def test_invalid_and_untrusted_metadata(self):
        for change in [ {"draft": True}, {"prerelease": True}, {"tag_name": "v1x2x3"} ]:
            with self.assertRaises(ValueError): mirror.release_assets({**self.release(), **change})
        release = self.release(); release["assets"][0]["digest"] = ""
        with self.assertRaises(ValueError): mirror.release_assets(release)
        release = self.release(); release["assets"][0]["browser_download_url"] = "https://example.com/installer.exe"
        with self.assertRaises(ValueError): mirror.release_assets(release)
    def test_publish_verified_stable_names_and_preserve_other_platform(self):
        client = FakeR2()
        client.objects["easyai/latest.json"] = json.dumps({"version":"v1.2.2","platforms":{"macos-x64":{"version":"v1.2.2"}}}).encode()
        release = self.release(); items = mirror.release_assets(release)
        with tempfile.TemporaryDirectory() as directory:
            directory = Path(directory); (directory/items["windows-x64"]["name"]).write_bytes(b"test")
            result = mirror.publish(client, "bucket", release, items, directory)
        self.assertEqual(client.objects["easyai/EasyAI-Windows-x64-Setup.exe"], b"test")
        self.assertEqual(result["platforms"]["macos-x64"]["version"], "v1.2.2")
        self.assertEqual(client.copies[0]["CacheControl"], mirror.STABLE_CACHE)
    def test_corrupt_staging_never_replaces_stable(self):
        client = FakeR2(); client.corrupt = True
        release = self.release(); items = mirror.release_assets(release)
        with tempfile.TemporaryDirectory() as directory:
            directory = Path(directory); (directory/items["windows-x64"]["name"]).write_bytes(b"test")
            with self.assertRaises(ValueError): mirror.publish(client,"bucket",release,items,directory)
        self.assertEqual(client.copies, [])
        self.assertNotIn("easyai/latest.json",client.objects)
    def test_old_version_cannot_overwrite_new(self):
        client = FakeR2(); client.objects["easyai/latest.json"] = b'{"version":"v2.0.0"}'
        with self.assertRaises(ValueError): mirror.publish(client,"bucket",self.release(),mirror.release_assets(self.release()),Path("unused"))
        self.assertEqual(client.copies, [])
    def test_download_rejects_corrupted_cache(self):
        release = self.release(); item = mirror.release_assets(release)["windows-x64"]
        with tempfile.TemporaryDirectory() as directory:
            directory = Path(directory); (directory/item["name"]).write_bytes(b"test")
            self.assertEqual(mirror.download(item,directory).read_bytes(),b"test")

if __name__ == "__main__": unittest.main()
