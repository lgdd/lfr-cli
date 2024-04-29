---
layout: default
title: command
grand_parent: Commands
parent: create
nav_order: 5
permalink: /cmd/create/command
---

# lfr create command

It allows you to create an OSGi module to add a new Gogo Shell command to Liferay DXP/Portal.

## Usage:
```shell
lfr create command NAME [flags]
# or
lfr c cmd NAME [flags]
```

## Flags:
- `-h`, `--help`
  - help for `lfr create command`

## Global Flags:
- `--no-color`
  - disable colors for output messages
- `-p`, `--package`
  - base package name (default `org.acme`)