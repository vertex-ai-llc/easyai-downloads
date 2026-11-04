#!/usr/bin/env bash
set -euo pipefail
artifact_dir="${1:-dist}"
: "${VERSION:?VERSION is required}"
: "${SOURCE_SHA:?SOURCE_SHA is required}"
: "${GITHUB_REPOSITORY:?GITHUB_REPOSITORY is required}"
release_tag="v$VERSION"
release_notes="${RELEASE_NOTES:-EasyAI desktop release. Source commit: $SOURCE_SHA}"

retry() {
  local attempt
  for attempt in 1 2 3 4; do
    if "$@"; then return 0; fi
    if [[ $attempt == 4 ]]; then return 1; fi
    echo "GitHub publish request failed; retrying ($attempt/4)." >&2
    sleep "$((attempt * 5))"
  done
}

ensure_release() {
  # Recheck existence after an ambiguous failure: creation may have succeeded
  # on GitHub even if the client received a 502 or lost the response.
  if gh release view "$release_tag" --repo "$GITHUB_REPOSITORY" >/dev/null 2>&1; then return 0; fi
  gh release create "$release_tag" --repo "$GITHUB_REPOSITORY" --target main --draft --title "EasyAI $VERSION" --notes "$release_notes"
}

test -s "$artifact_dir/SHA256SUMS"
retry ensure_release
shopt -s nullglob
assets=("$artifact_dir"/*)
for asset in "${assets[@]}"; do
  [[ -f "$asset" ]] || continue
  retry gh release upload "$release_tag" --repo "$GITHUB_REPOSITORY" --clobber "$asset"
done
# A newly created release remains a draft if any upload fails permanently.
retry gh release edit "$release_tag" --repo "$GITHUB_REPOSITORY" --title "EasyAI $VERSION" --notes "$release_notes" --draft=false --latest
