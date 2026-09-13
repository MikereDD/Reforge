#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

binary="$root_dir/bin/reforge"
catalog="$root_dir/manifests/catalog.example.json"
pacman_key="${REFORGE_PACMAN_KEY:-pacman-key}"
trust_db="${REFORGE_ARCH_TRUST_DB:-/etc/pacman.d/gnupg}"
keyring_source="${REFORGE_ARCH_KEYRING_SOURCE:-/usr/share/pacman/keyrings/archlinux.gpg}"
work_dir="${REFORGE_WORK_DIR:-$HOME/.cache/reforge/artifacts}"

if [[ ! -x "$binary" ]]; then
    printf 'Reforge binary is not executable: %s\n' "$binary" >&2
    printf 'Run: chmod +x "%s"\n' "$binary" >&2
    exit 1
fi

exec "$binary" \
    -catalog "$catalog" \
    -verify \
    -pacman-key "$pacman_key" \
    -arch-trust-db "$trust_db" \
    -arch-keyring-source "$keyring_source" \
    -work-dir "$work_dir" \
    "$@"
