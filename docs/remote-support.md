# Remote Support

## Goal

Reforge Remote Support allows a trusted technician to assist a user from the booted Reforge environment before the machine is repaired or reinstalled.

Remote support is not intended to become a permanent remote-access agent on the installed operating system.

## Session Model

1. User boots Reforge Key.
2. Reforge establishes networking.
3. User chooses **Remote Support**.
4. Reforge creates a short-lived session identity and displays a session code.
5. The technician connects using that code through the Reforge Hub/session broker.
6. Both sides authenticate the session.
7. Session credentials expire when the user ends support, reboots, or shuts down.

## Technician Capabilities

Initial read-oriented capabilities may include:

- hardware inventory
- storage layout and health
- mounted-volume inspection
- operating-system detection
- boot/EFI inspection
- logs and diagnostic reports
- network diagnostics

Later controlled actions may include:

- start a backup
- mount/unmount filesystems
- copy recovery files
- repair selected boot configuration
- prepare an installer
- request a destructive reinstall action

## Destructive Action Rule

Remote technicians must not be able to silently erase disks or begin equivalent destructive operations.

A destructive request must identify the exact target and be confirmed locally on the target machine.

Example:

```text
REMOTE REQUEST

Technician requests permission to erase:
Samsung SSD 990 PRO 2TB
Serial ending: 7A3F

ALL DATA ON THIS DEVICE WILL BE LOST.

[ DENY ]                 [ HOLD ENTER TO APPROVE ]
```

## Audit Trail

The rescue session should maintain a local action log containing:

- connection time
- technician/session identity
- read actions
- requested write/destructive actions
- local approval/denial
- execution result

Sensitive content such as file contents or secrets should not be indiscriminately logged.
