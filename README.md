# EasyAI Downloads

This public repository contains only the reproducible GitHub Actions builders
for EasyAI Desktop. Source code stays in a private repository and is checked
out by Actions with the `BUILD_SOURCE_TOKEN` secret.

## Release flow

1. A validated `v*` tag is pushed to the private source repository.
2. Source checks run there and dispatch `easyai-source-tag` to this repository.
3. This workflow checks out the exact immutable source tag, builds Windows and
   macOS artifacts, and publishes a GitHub Release with checksums.

The public repository does not contain source archives, credentials, signing
keys, or deployment configuration.

## Required configuration

Repository variable:

- `BUILD_SOURCE_REPOSITORY`: normally `vertex-ai-llc/easyai`

Repository secret:

- `BUILD_SOURCE_TOKEN`: a fine-grained token that can read the private source
  repository and has no write permission there.

Optional signing secrets can be added later to the Windows and macOS jobs.
