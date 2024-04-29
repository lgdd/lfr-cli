---
layout: default
title: diagnose
parent: Commands
nav_order: 1
permalink: /cmd/diagnose
---

# lfr diagnose

Inspired by the command `flutter doctor` from the [Flutter command-line tool](https://docs.flutter.dev/reference/flutter-cli#flutter-commands), it allows you to check your system to:
- See if Java is installed and which version is currently used.
- See if [Blade](https://learn.liferay.com/w/dxp/liferay-development/tooling/blade-cli/installing-and-updating-blade-cli){:target="_blank"} is installed.
- See if [Docker](https://docs.docker.com/get-docker/) is installed and report how much space is taken by official Liferay Docker Images and official Elasticsearch Docker Images.
- Report how much space is taken by the Liferay bundles stored under `~/.liferay/bundles` by Gradle or Maven when you intialize a Liferay Workspace.

## Usage:
```shell
lfr diagnose [flags]
# or
lfr diag
```

## Result example:
```shell
[✓] Java intalled (11.0.23)
    • Make sure that your Java edition is a Java Technical Compatibility Kit (TCK) compliant build.
    • JDK compatibility is for runtime and project compile time. DXP source compile is compatible with JDK 8 only.
[✗] Blade is not installed.
    • You might like this tool, but Blade is still the official one with useful features.
    • Blade is supported by Liferay and used by Liferay IDE behind the scenes.
    • Checkout the documentation: https://learn.liferay.com/w/dxp/building-applications/tooling/blade-cli
[✓] Docker installed (Docker Desktop 4.29.0 (145265))

[!] Downloaded bundles are using ~2.0 GB.
    • They are stored under ~/.liferay/bundles
[!] Official Liferay Docker images are using ~4.6 GB.
    • Run 'docker images liferay/dxp' to list DXP Images (EE)
    • Run 'docker images liferay/portal' to list Portal Images (CE)

More information about compatibilities: https://www.liferay.com/compatibility-matrix
```

## Flags:
- `-h`, `--help`
  - help for `lfr diagnose`

## Global Flags:
- `--no-color`
  - disable colors for output messages