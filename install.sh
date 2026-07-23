#!/usr/bin/env bash
# Builds Slyrics from source and installs it to ~/.local/bin (override with
# BINDIR). Run from anywhere inside the repo:
#
#   ./install.sh
set -e

repo_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bindir="${BINDIR:-$HOME/.local/bin}"

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required but wasn't found on your PATH." >&2
  echo "Install it from https://go.dev/dl/ and try again." >&2
  exit 1
fi

echo "Building slyrics-bin..."
( cd "$repo_dir" && go build -ldflags '-w -s' -o slyrics-bin . )

mkdir -p "$bindir"
install -m 755 "$repo_dir/slyrics-bin" "$bindir/slyrics-bin"
install -m 755 "$repo_dir/scripts/slyrics" "$bindir/slyrics"
rm -f "$repo_dir/slyrics-bin"

echo "Installed slyrics and slyrics-bin to $bindir"

case ":$PATH:" in
  *":$bindir:"*) ;;
  *)
    echo
    echo "$bindir is not on your PATH. Add this to your shell config:"
    echo "  export PATH=\"$bindir:\$PATH\""
    ;;
esac
