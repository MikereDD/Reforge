# Architecture

## Overview

Reforge is split into a lightweight boot client and remote services rather than being built around a permanent collection of ISO images.

```text
                    Official upstream sources
                 Microsoft / Arch / distro mirrors
                              |
                              v
+-------------+       +-------------------+
| Reforge Key | <---- |   Reforge Hub     |
|-------------|       |-------------------|
| Bootstrap   |       | Signed catalog    |
| Networking  |       | Profiles          |
| Rescue      |       | Session broker    |
| Diagnostics |       | Policy            |
| Cache       |       | Optional cache    |
+-------------+       +-------------------+
       |
       v
 Target machine
```

## Reforge Key

The Reforge Key is the bootable USB component.

Responsibilities:

- UEFI boot and bootstrap
- wired and wireless networking where supported
- hardware discovery
- disk and filesystem inspection
- one dependable offline rescue environment
- secure retrieval of the Reforge catalog
- payload download/stream/bootstrap orchestration
- temporary cache/workspace
- remote-support client

The key should not depend on storing every supported operating-system image locally.

## Reforge Hub

The Hub is primarily a control plane, not a mandatory mirror of every installer payload.

Responsibilities:

- publish signed installer/tool catalogs
- publish supported release/channel policy
- serve profile metadata and automation assets
- coordinate temporary remote-support sessions
- provide optional local or remote payload caching
- expose version/compatibility metadata

Where practical, payload bytes should come directly from official upstream infrastructure.

## Installer Catalog

The catalog is the authoritative description of what Reforge is willing to offer.

A catalog entry can define:

- stable identifier and display name
- menu order
- operating-system family
- supported architecture
- channel/release policy
- source policy
- boot method
- payload resolver or URL metadata
- hash/signature requirements
- tested/recommended state
- online/offline availability

The user-facing menu should be generated from this catalog rather than hard-coded around ISO filenames.

## Initial Installer Order

1. Windows 11
2. Arch Linux
3. Linux Mint
4. Ubuntu LTS
5. Fedora Workstation
6. Debian

## Reforge Rescue

Reforge Rescue remains available when remote installation is impossible.

Responsibilities:

- inspect SMART/NVMe health
- partition and filesystem inspection
- filesystem repair helpers
- mount/copy data
- backup/restore helpers
- networking diagnostics
- boot configuration inspection
- diagnostic report generation

## Remote Support

Remote support is initiated from the Reforge boot environment and uses an ephemeral session identity.

A technician may inspect the machine and request actions, but destructive requests require physical confirmation at the target machine.

See `docs/remote-support.md` and `docs/security-model.md`.

## Profiles

Profiles describe repeatable recovery/bootstrap preferences for a machine or person.

Potential profile content:

- preferred operating system
- hostname
- backup selection
- application inventory
- package selections
- driver references
- desktop environment
- filesystem preference
- post-install scripts
- restore notes

Secrets should not be stored directly in ordinary profile files.

## Optional Ventoy Fallback

Ventoy may remain as an optional compatibility/fallback path for manually supplied or unusual ISO images, but it is no longer the architectural foundation of Reforge.
