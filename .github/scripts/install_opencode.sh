#!/usr/bin/env bash
# Install the opencode CLI at the version pinned in ci/package.json
# and put its bin directory on the PATH.
# Needs: GITHUB_WORKSPACE, GITHUB_PATH.
set -euo pipefail
npm ci --prefix ci --ignore-scripts=false
echo "$GITHUB_WORKSPACE/ci/node_modules/.bin" >> "$GITHUB_PATH"
ci/node_modules/.bin/opencode --version
