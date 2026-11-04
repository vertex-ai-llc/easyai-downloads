"""Recover only successful platform builds from an immutable source version."""
import hashlib
import io
import json
import os
from pathlib import Path
import re
import subprocess
import urllib.error
import urllib.request
import zipfile


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


repo = os.environ['GITHUB_REPOSITORY']
token = os.environ['GH_TOKEN']
source_ref = os.environ['SOURCE_REF']
run_id = os.environ['BUILD_RUN_ID']
require(re.fullmatch(r'v[0-9]+\.[0-9]+\.[0-9]+', source_ref), 'Expected a release version tag')
require(run_id.isdigit(), 'Invalid build run ID')
version = json.loads(Path('source/apps/desktop/package.json').read_text())['version']
sha = subprocess.check_output(['git', '-C', 'source', 'rev-parse', 'HEAD'], text=True).strip()
tag_sha = subprocess.check_output(['git', '-C', 'source', 'rev-list', '-n', '1', source_ref], text=True).strip()
require(source_ref == 'v' + version and sha == tag_sha, 'Source version/tag mismatch')
api = 'https://api.github.com/repos/' + repo
headers = {'Authorization': 'Bearer ' + token, 'Accept': 'application/vnd.github+json'}


def get(path):
    with urllib.request.urlopen(urllib.request.Request(api + path, headers=headers), timeout=60) as response:
        return json.load(response)


def archive(artifact):
    request = urllib.request.Request(api + '/actions/artifacts/' + str(artifact['id']) + '/zip', headers=headers)
    try:
        response = urllib.request.build_opener(NoRedirect).open(request, timeout=60)
    except urllib.error.HTTPError as error:
        require(error.code in (301, 302, 303, 307, 308), 'Artifact download failed')
        target = error.headers['Location']
        require(target.startswith('https://'), 'Artifact redirect must use HTTPS')
        # Never send the GitHub bearer credential to artifact storage.
        response = urllib.request.urlopen(target, timeout=120)
    with response:
        content = response.read()
    expected = artifact.get('digest', '')
    require(expected.startswith('sha256:'), 'Missing artifact archive digest')
    require(hashlib.sha256(content).hexdigest() == expected[7:], 'Artifact archive digest mismatch')
    return zipfile.ZipFile(io.BytesIO(content))


run = get('/actions/runs/' + run_id)
require(run['path'] == '.github/workflows/build-release.yml', 'Artifacts must come from the release build workflow')
jobs = []
for attempt in range(1, run['run_attempt'] + 1):
    jobs.extend(get(f'/actions/runs/{run_id}/attempts/{attempt}/jobs?per_page=100')['jobs'])
artifacts = get('/actions/runs/' + run_id + '/artifacts?per_page=100')['artifacts']
dist = Path('dist')
dist.mkdir(exist_ok=True)
require(not list(dist.iterdir()), 'Recovery output directory must be empty')
platforms = []
for platform, filename, job_matches in [
    ('Windows-x64', f'EasyAI-{version}-Windows-x64-Setup.exe', lambda name: name == 'windows'),
    ('macOS-arm64', f'EasyAI-{version}-macOS-arm64.dmg', lambda name: name.startswith('macos (') and ', arm64)' in name),
    ('macOS-x64', f'EasyAI-{version}-macOS-x64.dmg', lambda name: name.startswith('macos (') and ', x64)' in name),
]:
    name = f'EasyAI-{platform}-{version}-{sha}'
    matching = [a for a in artifacts if a['name'] == name and not a['expired']]
    if not matching:
        continue
    require(any(j['conclusion'] == 'success' and job_matches(j['name']) for j in jobs), 'Platform build did not pass: ' + platform)
    require(len(matching) == 1, 'Ambiguous platform artifact: ' + platform)
    with archive(matching[0]) as package:
        entries = {item.filename: item for item in package.infolist() if not item.is_dir()}
        require(filename in entries, 'Missing installer in archive')
        data = package.read(filename)
        require(len(data) > 100_000, 'Installer is unexpectedly small')
        (dist / filename).write_bytes(data)
        if platform == 'Windows-x64':
            metadata = json.loads(package.read('release.json').decode('utf-8-sig'))
            require(metadata['version'] == version and metadata['artifact'] == filename and metadata['codeSigned'] is True, 'Windows signing metadata mismatch')
            require(hashlib.sha256(data).hexdigest() == metadata['sha256'], 'Windows installer digest mismatch')
            (dist / 'release.json').write_text(json.dumps(metadata, indent=2) + '\n')
    platforms.append(platform)
require(platforms, 'No successfully verified platform artifacts')
try:
    existing = get('/releases/tags/' + source_ref)
except urllib.error.HTTPError as error:
    if error.code != 404:
        raise
else:
    require(not existing['assets'], 'Use the normal complete release workflow to update an existing release with assets')
(dist / 'SHA256SUMS').write_text(''.join(hashlib.sha256(p.read_bytes()).hexdigest() + '  ' + p.name + '\n' for p in sorted(dist.iterdir()) if p.is_file()))
missing = sorted(set(['Windows-x64', 'macOS-arm64', 'macOS-x64']) - set(platforms))
notes = f'EasyAI desktop release. Source commit: {sha}\n\nAvailable: ' + ', '.join(platforms) + '.'
if missing:
    notes += '\nPending build: ' + ', '.join(missing) + '. These packages are not included in this release yet.'
with open(os.environ['GITHUB_ENV'], 'a') as env:
    env.write(f'VERSION={version}\nSOURCE_SHA={sha}\n')
    env.write('RELEASE_NOTES<<EASYAI_RELEASE_NOTES\n' + notes + '\nEASYAI_RELEASE_NOTES\n')
print(json.dumps({'source': sha, 'version': version, 'build_run_id': run_id, 'verified_platforms': platforms, 'pending_platforms': missing}))
