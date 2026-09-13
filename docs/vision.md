# Vision

Reforge exists to make recovery and reinstallation deliberate instead of impulsive.

A machine should not need a thumb drive containing dozens of aging ISO images. Reforge should instead provide a small, dependable boot key that can diagnose the machine, preserve data, reach approved remote installation sources, and help rebuild the system in a repeatable way.

## Product Direction

Reforge is an **internet-first, remote-capable recovery and deployment platform** delivered through a bootable USB key.

The key should remain useful when the internet is unavailable, but network access unlocks the full installer catalog, remote assistance, profiles, current packages, and optional cached payloads.

## Core Experiences

- guided reinstall preparation
- hardware and storage diagnostics
- backup before destructive work
- network-first Windows and Linux installation
- offline rescue and repair
- machine/profile-aware restore
- post-install bootstrap
- temporary remote-support sessions
- diagnostic reports that can recommend repair instead of reinstall

## Guiding Principles

1. **Safety first** — preserve data and identify the correct target before destructive actions.
2. **Network first, not network only** — use remote sources for large payloads while retaining meaningful offline rescue capability.
3. **Small key, large capability** — the USB carries the bootstrap and rescue environment, not an ISO warehouse.
4. **Official sources first** — retrieve payloads from Microsoft and Linux distribution infrastructure whenever practical.
5. **Verify before booting or installing** — catalog metadata and cryptographic verification are part of the installation path.
6. **Local consent for destructive remote actions** — remote support may prepare operations, but erasure/install requests require physical approval.
7. **Ephemeral remote access** — Reforge support sessions should expire when the rescue session ends unless the user explicitly installs something persistent.
8. **Clear over clever** — guided workflows should remain understandable even when the implementation is sophisticated.
9. **Extensible** — additional operating systems, tools, boot methods, and profiles should be catalog-driven where possible.
10. **Reproducible** — the project should be rebuildable from source, manifests, and documented tooling.
