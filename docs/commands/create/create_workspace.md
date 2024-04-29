---
layout: default
title: lfr create workspace
grand_parent: Commands
parent: lfr create
nav_order: 1
permalink: /commands/create/workspace
---

# lfr create workspace

It allows you to create a Liferay Workspace and specify if you want to use Gradle or Maven, DXP (EE) or Portal (CE) and on which major version.

Without flags, it creates a Gradle workspace on the latest version and release of Liferay Portal (CE). You can change defaults (i.e. without explicit flags) by editing the `~/.lfr/config.toml` file.

For DXP releases, the Liferay Workspace will have an custom URL for the bundle (`liferay.workspace.bundle.url` commented by default in the `gradle.properties` or `pom.xml`). Because the `releases-cdn.liferay.com` URLs are often slow or unresponsive, this [custom repository](https://github.com/lgdd/liferay-dxp-releases){:target="_blank"} is mirroring the Tomcat bundles of each release as a GitHub release. It should improve the `lfr init` command.

## Usage:
```shell
lfr create workspace NAME [flags]
# or
lfr c ws NAME [flags]
```

## Examples:
```shell
lfr c ws my-workspace
lfr c ws my-workspace -e dxp
lfr c ws my-workspace -e dxp -b maven
lfr c ws my-workspace -e dxp -b maven -v 7.3
```

## Flags:
- `-h`, `--help`
  - help for `lfr create`

## Global Flags:
- `--no-color`
  - disable colors for output messages
- `-p`, `--package string`
  - base package name (default "org.acme")