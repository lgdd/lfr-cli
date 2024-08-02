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

{: .new }
> Starting with **v3.1.0**, the command also check if [LCP](https://learn.liferay.com/w/liferay-cloud/reference/command-line-tool) is installed.

## Usage:
```shell
lfr diagnose [flags]
# or
lfr diag
```

## Result example:
```shell
[!] Java intalled (11.0.24)
    ! Liferay DXP DXP 2024.Q2 and Liferay Portal 7.4 GA120 will be the last version to support Java 11.
    • Make sure that your Java edition is a Java Technical Compatibility Kit (TCK) compliant build.
    • JDK compatibility is for runtime and project compile time.
[✗] Blade is not installed.
    • You might like this tool, but Blade is still the official one with useful features.
    • Blade is supported by Liferay and used by Liferay IDE behind the scenes.
    • Checkout the documentation: https://learn.liferay.com/w/dxp/building-applications/tooling/blade-cli
[!] LCP is not installed.
    • If you work on Liferay PaaS or Liferay SaaS, LCP can be used to view and manage your Liferay Cloud services.
    • Checkout the documentation: https://learn.liferay.com/w/liferay-cloud/reference/command-line-tool
[✓] Docker installed (27.1.1)

[!] Downloaded bundles are using ~1.9 GB.
    • They are stored under /home/lgd/.liferay/bundles
[!] Official Liferay Docker images are using ~2.2 GB.
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