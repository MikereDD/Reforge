# Security Model

## Trust Boundaries

Reforge interacts with several distinct trust domains:

- the local Reforge Key
- the target machine and its disks
- the Reforge Hub
- official upstream distribution/vendor infrastructure
- optional technician clients
- local-network caches

## Core Rules

### Verify catalogs

Installer catalogs must be authenticated before Reforge trusts URLs, boot parameters, hashes, or policy metadata.

### Verify payloads

Downloaded payloads must be cryptographically verified according to the catalog entry before they are executed, mounted, or used for installation.

### Prefer HTTPS and official sources

Official HTTPS sources are preferred. Mirrors and fallback sources must still satisfy integrity verification.

### No silent destructive remote control

Remote support may inspect and prepare operations, but destructive actions require explicit local approval.

### Ephemeral support identity

Support-session credentials are temporary and should not survive reboot/shutdown unless a future feature explicitly and separately enables persistence.

### Minimal secret storage

Ordinary Git-tracked profiles and catalogs must not contain passwords, API tokens, recovery keys, or private signing keys.

### Distinguish data disks clearly

Disk-selection UI must show enough identifying information to reduce accidental erasure, including size, model, serial suffix, partitions, and detected operating systems when available.

### Fail closed on integrity errors

If a required signature/hash cannot be validated, Reforge should not continue with that payload unless an explicit advanced override policy is later designed.

## Future Work

- catalog signing format and key rotation
- Secure Boot chain
- Hub authentication
- technician identity and authorization
- encrypted profile storage
- audit-log integrity
- supply-chain and mirror-fallback policy
