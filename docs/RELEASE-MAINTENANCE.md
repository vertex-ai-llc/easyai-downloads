# EasyAI Downloads

Build automation and published desktop releases.

Application source and release credentials are kept outside this public
repository. Unsigned builds are for testing; signing can be enabled by the
maintainer for production releases.

If one hosted runner is unavailable, run **Build EasyAI Desktop** with the
immutable source tag and set verified_build_run_id to the original build run ID. It reuses only
artifacts whose platform job succeeded, matches the exact source SHA, verifies
the archive digest and Windows signing metadata, and lists missing platforms
in the release notes. It never rebuilds or substitutes another version.
This recovery entry point creates a release only when it has no existing assets;
use the normal full build/publish workflow to update an existing release.

Run python scripts/test-recover-verified-artifacts.py to check recovery gates.

## Stable R2 downloads

After **Build EasyAI Desktop** succeeds (normal or verified recovery),
**Refresh download links** automatically mirrors the latest published release
to R2, then updates the homepage. It also supports manual dispatch without
rebuilding the app. Publication is serialized; failures keep the existing
homepage and are reported as a failed workflow.

Public base: `https://easyai-download.postlyn.com`

| Platform | Fixed object key |
| --- | --- |
| Windows x64 | `easyai/EasyAI-Windows-x64-Setup.exe` |
| macOS Apple Silicon | `easyai/EasyAI-macOS-arm64.dmg` |
| macOS Intel | `easyai/EasyAI-macOS-x64.dmg` |
| Version and SHA-256 manifest | `easyai/latest.json` |

Installers are checked against GitHub's asset size and SHA-256, staged at
`easyai/releases/<tag>/`, fully read back and hashed, then copied to the fixed
names and checked again. The manifest is written last. A partial release keeps
the previously working objects for other platforms; the homepage only promotes
an R2 link whose version, size and digest match the current release. Older
versions cannot replace a newer manifest. Fixed names require revalidation
instead of retaining stale downloads in a long-lived cache.

Repository Actions secret **R2_UPLOAD_CONFIG** is a JSON object with `bucket`,
`endpoint`, `access_key_id`, and `secret_access_key`. Use an R2 token scoped to
the download bucket. Never commit this JSON. The updater requires Python 3.11+
and boto3. Local run: `python scripts/mirror-release-r2.py --config <private-file>`.
Tests: `python scripts/test-mirror-release-r2.py` and
`python scripts/test-update-homepage.py`.

## Official ChatGPT and Claude installers

Run **Refresh official agent installers** when the official packages need
refreshing. Select all, Claude, or ChatGPT. It runs the existing EasyAI
maintenance tool from the selected source tag using the same R2 secret;
`BUILD_SOURCE_REPOSITORY` and `BUILD_SOURCE_TOKEN` are the existing source
checkout settings. This is independent of EasyAI release publishing.

The fixed keys remain `agents/chatgpt/ChatGPT-x64.msix`,
`agents/chatgpt/ChatGPT_Classic.dmg`, `agents/claude/Claude-Latest.exe`,
and `agents/claude/Claud-Latest.dmg` (the existing spelling is retained).
The updater reads official latest redirects, validates installer formats,
uploads complete packages, and verifies stored source metadata and sizes.
If an official endpoint blocks a download, the workflow fails visibly and
preserves that package's previous object. A metadata artifact records the
resolved official source URLs.
