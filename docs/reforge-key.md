# Reforge Key Development Bundle

Reforge Key is the small Linux-side runtime that resolves, verifies, and eventually boots approved recovery and installer targets.

The current development bundle is not yet a bootable USB image. It exists to validate the Reforge runtime on Linux before the boot-image layer is added.

## Arch trust model

Reforge verifies Arch netboot artifacts through the host's initialized pacman trust database.

The packaged keyring file:

```text
/usr/share/pacman/keyrings/archlinux.gpg
```

is treated as evidence that the `archlinux-keyring` package is present. It is **not** passed directly to `gpgv`.

Actual detached-signature verification uses:

```text
pacman-key -v <signature> <artifact>
```

which consults pacman's operational GnuPG trust database beneath:

```text
/etc/pacman.d/gnupg
```

Reforge does not download signing keys on demand.

## Trust requirements

The Linux Reforge Key environment requires:

- `pacman-key`
- an initialized `/etc/pacman.d/gnupg`
- `/usr/share/pacman/keyrings/archlinux.gpg`
- a usable system CA certificate bundle

The preparation script checks these requirements and can install the required Arch packages when explicitly requested.

## Build from Windows

From the repository root:

```powershell
.\scripts\Build-Reforge-Key.ps1 -Architecture amd64 -Clean
```

For an ARM64 Linux test host:

```powershell
.\scripts\Build-Reforge-Key.ps1 -Architecture arm64 -Clean
```

Bundles are written beneath `dist/`, which is intentionally ignored by Git.

## Prepare the Linux host

After copying the bundle to an Arch Linux machine:

```bash
cd reforge-key-linux-<arch>
chmod +x bin/reforge run.sh scripts/prepare.sh
./scripts/prepare.sh
```

If the required packages are missing:

```bash
./scripts/prepare.sh --install
```

## Run a real verification

When preparation reports `Status: READY`:

```bash
./run.sh
```

Select Arch Linux and confirm the operation.

The intended progression is:

```text
Status: RESOLVED

Verification Preflight

pacman-key:        PASS
Pacman trust DB:   PASS
Arch keyring:      PASS

Status: READY

Artifact Verification

UEFI / x86_64
  Download:       PASS
  SHA-256:        ...
  PGP signature:  VALID

BIOS / x86_64
  Download:       PASS
  SHA-256:        ...
  PGP signature:  VALID

Status: VERIFIED
```

Any missing trust environment or invalid signature must leave the target blocked.
