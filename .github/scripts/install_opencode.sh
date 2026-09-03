#!/usr/bin/env bash
# Install the opencode CLI at the version pinned in .github/package.json
# and put its bin directory on the PATH.
# Needs: GITHUB_WORKSPACE, GITHUB_PATH.
set -euo pipefail
npm ci --prefix .github --ignore-scripts=false
echo "$GITHUB_WORKSPACE/.github/node_modules/.bin" >> "$GITHUB_PATH"
.github/node_modules/.bin/opencode --version
