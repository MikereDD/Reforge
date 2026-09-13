# Reforge

**Multi-OS Remote Recovery & Reinstall Toolkit**

Reforge is an internet-first recovery and deployment platform delivered through a small bootable USB key.

The USB is intentionally lightweight. It provides the Reforge bootstrap environment, networking, diagnostics, one dependable offline rescue path, and enough local tooling to reach approved remote installation and recovery resources. Large operating-system images do not need to live permanently on the thumb drive.

## Core Components

- **Reforge Key** — small bootable USB environment used to start recovery, diagnostics, installation, and remote-support workflows.
- **Reforge Hub** — remote control plane for signed catalogs, profiles, remote-support coordination, policy, and optional payload caching.
- **Reforge Rescue** — offline-capable rescue environment for storage, filesystem, backup, and network troubleshooting.
- **Installer Catalog** — approved operating-system and tool definitions with source, architecture, verification, and boot metadata.
- **Profiles** — reusable machine/user setup definitions for backup, restore, and post-install bootstrap.

## Initial Installer Order

1. Windows 11
2. Arch Linux
3. Linux Mint
4. Ubuntu LTS
5. Fedora Workstation
6. Debian

## Design Goals

- Keep the USB small and replaceable instead of storing a library of ISOs.
- Retrieve installation payloads from official upstream sources on demand whenever practical.
- Verify downloaded content before use.
- Preserve one useful offline rescue environment.
- Support remote assistance without installing a permanent backdoor.
- Require local confirmation for destructive remote actions.
- Make repeated reinstalls safer, faster, and more reproducible.

## Repository Layout

- `docs/` — project planning, architecture, security, and workflow documentation
- `rescue/` — Reforge Rescue assets and build work
- `scripts/` — Linux and Windows helpers
- `ventoy/` — optional/fallback Ventoy configuration and theme assets
- `profiles/` — reusable restore/install profile examples
- `manifests/` — installer catalog, schemas, checksums, and release metadata
- `tools/` — maintenance and build utilities
- `assets/` — branding assets

## Current Status

`v0.2-dev — Remote Installer Foundation`

The current development focus is the remote-first architecture and the catalog/verification system that will drive network installation.

## Repository Policy

Large installer images, downloaded binaries, backups, caches, and generated rescue images should not be committed to Git. Reforge source control should contain the logic and metadata needed to reproduce the toolkit, not the payload library itself.
