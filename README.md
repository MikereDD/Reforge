# Reforge

**Multi-OS Recovery & Reinstall Toolkit**

Reforge is a Ventoy-based recovery and reinstall toolkit that combines:

- multi-OS installer hosting
- a customized rescue environment
- pre-reinstall backup helpers
- post-install restore and bootstrap helpers
- reusable machine/user profiles
- manifests and integrity verification

## Initial Goals

- Boot multiple Linux installers and Windows from one USB device.
- Provide a guided rescue environment before destructive reinstalls.
- Preserve important user files, driver exports, and app inventories.
- Restore a machine faster after a fresh install.
- Keep the toolkit maintainable with manifests, documentation, and versioned scripts.

## Planned Structure

- `docs/` — project planning, roadmap, architecture, workflow docs
- `rescue/` — custom rescue-environment assets and build notes
- `scripts/` — Linux and Windows helper scripts
- `ventoy/` — Ventoy theme/config assets
- `profiles/` — reusable restore/install profile examples
- `manifests/` — checksums, inventories, release metadata
- `tools/` — helper utilities and maintenance scripts
- `assets/` — branding assets such as icons/logos

## Current Status

This repository is the initial project foundation.

## Notes

- ISO images, downloaded binaries, large backups, and generated rescue images should **not** be committed to Git.
- Keep the repository focused on build/configuration logic, documentation, and automation.
