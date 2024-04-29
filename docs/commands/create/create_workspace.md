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

The command is also running a `git init` for the created workspace, and is adding a [GitHub Action](https://github.com/lgdd/lfr-cli/blob/main/internal/assets/tpl/github/liferay-upgrade.yml){:target="_blank"} to help you upgrade the Liferay Workspaces you host on GitHub.

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

For DXP releases, the Liferay Workspace will find a custom bundle URL in the `gradle.properties` or `pom.xml` (`liferay.workspace.bundle.url` commented by default ). Because the `releases-cdn.liferay.com` URLs are often slow or unresponsive, this [custom repository](https://github.com/lgdd/liferay-dxp-releases){:target="_blank"} is mirroring the Tomcat bundles of each release as a GitHub release. It should improve the `lfr init` command.

## Flags:
- `-h`, `--help`
  - help for `lfr create`

## Global Flags:
- `--no-color`
  - disable colors for output messages
- `-p`, `--package string`
  - base package name (default "org.acme")