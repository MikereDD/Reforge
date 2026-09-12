# Architecture

## High-Level Layout

- **Ventoy layer**
  - boots installer and rescue ISOs
  - provides consistent USB boot experience

- **Reforge Rescue layer**
  - customized rescue environment
  - guided TUI/CLI workflows
  - backup, diagnostics, and repair helpers

- **Profile layer**
  - user/machine-specific restore manifests
  - optional application lists, driver capture, and restore hints

- **Automation layer**
  - Windows PowerShell helpers
  - Linux shell helpers
  - download/update/verification scripts
