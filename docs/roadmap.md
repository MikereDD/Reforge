# Roadmap

## v0.1 — Foundation — Complete

- Establish repository structure.
- Add baseline documentation.
- Add placeholder Linux and Windows scripts.
- Add initial profile/manifests conventions.
- Add Reforge branding and icon.

## v0.2-dev — Remote Installer Foundation

- Replace the ISO-library model with Reforge Key + Reforge Hub architecture.
- Define the installer catalog schema.
- Add Windows 11 and Arch Linux as the first two installer targets.
- Define supported installer ordering and metadata.
- Add source and integrity-verification policy.
- Define remote-support architecture and security rules.
- Prototype catalog loading and validation.
- Prototype payload resolution without retaining permanent ISO copies.

## v0.3-dev — Reforge Key

- Build the first bootable Reforge Key image.
- Add network initialization and connectivity diagnostics.
- Add hardware/storage discovery.
- Bundle Reforge Rescue fallback capability.
- Add local workspace/cache management.
- Render the install/recovery menu from the catalog.

## v0.4-dev — Network Installation

- Implement Arch Linux network installation path.
- Implement Windows 11 remote-media/bootstrap path.
- Add Linux Mint, Ubuntu LTS, Fedora Workstation, and Debian paths.
- Verify all downloaded payloads before use.
- Add retry/resume and mirror/source fallback behavior.

## v0.5-dev — Rescue & Reinstall Preparation

- Add guided backup workflow.
- Add Windows installation inventory and driver export.
- Add Linux package/config inventory.
- Add filesystem, SMART/NVMe, and boot diagnostics.
- Add Reforge diagnostic report generation.

## v0.6-dev — Remote Support

- Add temporary support-session codes.
- Add authenticated technician connection.
- Add remote hardware/log/storage inspection.
- Add auditable action requests.
- Require physical confirmation for destructive operations.
- Ensure session credentials expire with the rescue session.

## v0.7-dev — Profiles & Restore

- Finalize profile schema.
- Add machine-aware profiles.
- Add application/package restore workflows.
- Add post-install bootstrap.
- Add optional Hub-side encrypted profile storage.

## v0.8-dev — Hub & Caching

- Implement optional Reforge Hub service.
- Add local-network payload cache.
- Add catalog signing and publication workflow.
- Add cache eviction/version policy.
- Add update/channel management.

## v0.9-dev — Hardening

- Secure Boot strategy.
- Threat-model review.
- Failure/recovery testing.
- Disk-selection safety testing.
- Offline behavior testing.
- Network interruption testing.
- Remote-support authorization review.

## v1.0 — Reforge

- Stable Reforge Key.
- Stable signed catalog.
- Supported Windows 11 and Linux installation workflows.
- Offline rescue capability.
- Backup/restore workflows.
- Remote-support workflow.
- Documented build/release process.
