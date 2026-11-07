#!/usr/bin/env sh
set -eu
cd "$(git rev-parse --show-toplevel)"
want="$(tr -d '[:space:]' < VERSION)"
out="$(mktemp -d)"
trap 'rm -rf "$out"' EXIT
go build -trimpath -buildvcs=false -ldflags="-X github.com/FacileStudio/Mycelium/cmd.version=${want}" -o "$out/mycelium" .
got="$("$out/mycelium" --version | awk '{print $2}')"
[ "$got" = "$want" ] || { echo "version mismatch: binary $got, VERSION $want" >&2; exit 1; }
top="$(sed -n 's/^## \[\([0-9][^]]*\)\].*/\1/p' CHANGELOG.md | head -1)"
[ "$top" = "$want" ] || [ "$top" = "Unreleased" ] || { echo "VERSION $want is not the top CHANGELOG release ($top)" >&2; exit 1; }
echo "version ok: $want"
