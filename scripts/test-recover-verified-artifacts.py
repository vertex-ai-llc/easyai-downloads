import hashlib
import io
import json
import os
from pathlib import Path
import runpy
import tempfile
import unittest
from unittest.mock import patch
import urllib.error
import zipfile


SCRIPT = Path(__file__).resolve().with_name('recover-verified-artifacts.py')
SHA = 'a' * 40
VERSION = '0.2.146'
FILENAME = f'EasyAI-{VERSION}-Windows-x64-Setup.exe'


class Response(io.BytesIO):
    pass


class RecoveryTests(unittest.TestCase):
    def run_fixture(self, *, corrupt=False, failed=False, wrong_source=False, existing=False):
        payload = b'installer-fixture' * 10000
        metadata = {'version': VERSION, 'artifact': FILENAME, 'codeSigned': True, 'sha256': hashlib.sha256(payload).hexdigest()}
        packed = io.BytesIO()
        with zipfile.ZipFile(packed, 'w') as archive:
            archive.writestr(FILENAME, payload)
            archive.writestr('release.json', json.dumps(metadata))
        data = packed.getvalue()
        artifact = {'id': 55, 'name': f'EasyAI-Windows-x64-{VERSION}-{SHA if not wrong_source else "b" * 40}', 'expired': False, 'digest': 'sha256:' + ('0' * 64 if corrupt else hashlib.sha256(data).hexdigest())}
        calls = []

        def urlopen(request, **kwargs):
            url = request.full_url if hasattr(request, 'full_url') else request
            calls.append(url)
            if url == 'https://storage.example/archive':
                self.assertIsInstance(request, str, 'Storage download must not carry GitHub headers')
                return Response(data)
            if url.endswith('/actions/runs/123'):
                value = {'path': '.github/workflows/build-release.yml', 'run_attempt': 2}
            elif '/attempts/' in url:
                value = {'jobs': [{'name': 'windows', 'conclusion': 'failure' if failed else 'success'}]}
            elif '/artifacts?' in url:
                value = {'artifacts': [artifact]}
            elif '/releases/tags/' in url:
                if not existing:
                    raise urllib.error.HTTPError(url, 404, 'Not found', {}, None)
                value = {'assets': [{'name': 'previous.exe'}]}
            else:
                raise AssertionError(url)
            return Response(json.dumps(value).encode())

        class Opener:
            def open(self, request, **kwargs):
                raise urllib.error.HTTPError(request.full_url, 302, 'Found', {'Location': 'https://storage.example/archive'}, None)

        cwd = Path.cwd()
        with tempfile.TemporaryDirectory() as temporary:
            os.chdir(temporary)
            try:
                package = Path('source/apps/desktop/package.json')
                package.parent.mkdir(parents=True)
                package.write_text(json.dumps({'version': VERSION}))
                env = {'GITHUB_REPOSITORY': 'fixture/repo', 'GH_TOKEN': 'fixture', 'SOURCE_REF': 'v' + VERSION, 'BUILD_RUN_ID': '123', 'GITHUB_ENV': str(Path(temporary) / 'environment')}
                with patch.dict(os.environ, env), patch('subprocess.check_output', return_value=SHA), patch('urllib.request.urlopen', side_effect=urlopen), patch('urllib.request.build_opener', return_value=Opener()):
                    runpy.run_path(str(SCRIPT), run_name='__main__')
                self.assertEqual(Path('dist', FILENAME).read_bytes(), payload)
                lines = Path('dist/SHA256SUMS').read_text().splitlines()
                self.assertEqual(len(lines), 2)
                for line in lines:
                    digest, name = line.split('  ')
                    self.assertEqual(digest, hashlib.sha256(Path('dist', name).read_bytes()).hexdigest())
                settings = Path(env['GITHUB_ENV']).read_text()
                self.assertIn('VERSION=0.2.146\nSOURCE_SHA=', settings)
                self.assertIn('Pending build: macOS-arm64, macOS-x64.', settings)
            finally:
                os.chdir(cwd)

    def test_verified_artifacts(self):
        self.run_fixture()

    def test_corrupt_archive_rejected(self):
        with self.assertRaisesRegex(RuntimeError, 'digest mismatch'):
            self.run_fixture(corrupt=True)

    def test_failed_build_rejected(self):
        with self.assertRaisesRegex(RuntimeError, 'did not pass'):
            self.run_fixture(failed=True)

    def test_wrong_source_rejected(self):
        with self.assertRaisesRegex(RuntimeError, 'No successfully verified'):
            self.run_fixture(wrong_source=True)

    def test_existing_assets_preserved(self):
        with self.assertRaisesRegex(RuntimeError, 'existing release with assets'):
            self.run_fixture(existing=True)


if __name__ == '__main__':
    unittest.main()

