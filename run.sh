#!/usr/bin/env bash
set -euo pipefail

bazel build //cli:cli --ui_event_filters=-info --noshow_progress 2>/dev/null

# .bazel directory is a symlink configured in .bazelrc
.bazel/bin/cli/cli_/cli "$@"
