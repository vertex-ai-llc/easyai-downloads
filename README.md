# EasyAI Downloads

Build automation and published desktop releases.

Application source and release credentials are kept outside this public
repository. Unsigned builds are for testing; signing can be enabled by the
maintainer for production releases.

If one hosted runner is unavailable, use **Publish verified EasyAI artifacts**
with the immutable source tag and the original build run ID. It reuses only
artifacts whose platform job succeeded, matches the exact source SHA, verifies
the archive digest and Windows signing metadata, and lists missing platforms
in the release notes. It never rebuilds or substitutes another version.
This recovery entry point creates a release only when it has no existing assets;
use the normal full build/publish workflow to update an existing release.

Run python scripts/test-recover-verified-artifacts.py to check recovery gates.
