import importlib.util
import tempfile
import unittest
from pathlib import Path
spec = importlib.util.spec_from_file_location("homepage", Path(__file__).with_name("update-homepage.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)

class DownloadLinks(unittest.TestCase):
    def release(self):
        name = "EasyAI-1.2.3-Windows-x64-Setup.exe"
        return {"tag_name": "v1.2.3", "published_at": "2026-11-05T18:00:00Z", "assets": [{"name": name, "size": 10, "browser_download_url": f"{module.BASE}/releases/download/v1.2.3/{name}"}]}
    def test_partial_release_uses_only_existing_assets(self):
        block = module.download_block(self.release())
        self.assertIn("Windows-x64-Setup.exe", block)
        self.assertNotIn(".dmg", block)
        self.assertIn("2026-11-06 发布", block)
        self.assertEqual(block.count("查看可用版本"), 2)
    def test_rejects_unpublished_and_foreign_urls(self):
        for key in ["draft", "prerelease"]:
            release = self.release(); release[key] = True
            with self.assertRaises(ValueError): module.download_block(release)
        release = self.release(); release["assets"][0]["browser_download_url"] = "https://example.com/file.exe"
        with self.assertRaises(ValueError): module.download_block(release)
    def test_preserves_copy_and_qr_and_is_idempotent(self):
        with tempfile.TemporaryDirectory() as folder:
            readme = Path(folder)/"README.md"
            readme.write_text("intro\n"+module.START+"old"+module.END+"\ncommunity-qr", encoding="utf-8")
            self.assertTrue(module.update(readme,self.release()))
            self.assertFalse(module.update(readme,self.release()))
            self.assertTrue(readme.read_text(encoding="utf-8").endswith("community-qr"))

if __name__ == "__main__": unittest.main()
