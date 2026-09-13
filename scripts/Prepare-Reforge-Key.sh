#!/usr/bin/env bash
set -euo pipefail

install_missing=false

usage() {
    cat <<'EOF'
Usage: Prepare-Reforge-Key.sh [--install]

Checks the Arch Linux trust environment required by Reforge.

Options:
  --install   Install missing Arch Linux packages with pacman.
  -h, --help  Show this help.
EOF
}

while (($# > 0)); do
    case "$1" in
        --install)
            install_missing=true
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            printf 'Unknown argument: %s\n' "$1" >&2
            usage >&2
            exit 2
            ;;
    esac
    shift
done

if [[ ! -e /etc/arch-release ]]; then
    printf 'BLOCKED: Reforge Key preparation currently expects Arch Linux.\n' >&2
    exit 1
fi

if ! command -v pacman >/dev/null 2>&1; then
    printf 'BLOCKED: pacman is unavailable.\n' >&2
    exit 1
fi

required_packages=(
    archlinux-keyring
    gnupg
    ca-certificates
)

if "$install_missing"; then
    if ((EUID == 0)); then
        pacman -S --needed "${required_packages[@]}"
    else
        if ! command -v sudo >/dev/null 2>&1; then
            printf 'BLOCKED: sudo is required for --install when not running as root.\n' >&2
            exit 1
        fi

        sudo pacman -S --needed "${required_packages[@]}"
    fi
fi

failed=false

printf 'Reforge Key trust preflight\n\n'

if command -v pacman-key >/dev/null 2>&1; then
    printf 'pacman-key:        PASS (%s)\n' "$(command -v pacman-key)"
else
    printf 'pacman-key:        MISSING\n'
    failed=true
fi

trust_db=/etc/pacman.d/gnupg

if [[ -d "$trust_db" ]]; then
    printf 'Pacman trust DB:   PASS (%s)\n' "$trust_db"
else
    printf 'Pacman trust DB:   MISSING (%s)\n' "$trust_db"
    failed=true
fi

keyring=/usr/share/pacman/keyrings/archlinux.gpg

if [[ -f "$keyring" ]]; then
    printf 'Arch keyring:      PASS (%s)\n' "$keyring"
else
    printf 'Arch keyring:      MISSING (%s)\n' "$keyring"
    failed=true
fi

ca_bundle=/etc/ssl/certs/ca-certificates.crt

if [[ -f "$ca_bundle" ]]; then
    printf 'CA bundle:         PASS (%s)\n' "$ca_bundle"
else
    printf 'CA bundle:         MISSING (%s)\n' "$ca_bundle"
    failed=true
fi

printf '\n'

if "$failed"; then
    printf 'Status: BLOCKED\n'

    if ! "$install_missing"; then
        printf '\nRun again with --install to install the required Arch packages.\n'
    fi

    exit 1
fi

printf 'Status: READY\n'
