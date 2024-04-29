---
layout: default
title: start
parent: Commands
nav_order: 6
permalink: /cmd/start
---

# lfr start

It allows you to start the Tomcat bundle created from the `lfr init` command. When starting the bundle, it saves a `liferay.pid` file at the root of the Liferay Workspace to allow to stop or kill the process (e.g. `kill (cat liferay.pid)`).

## Usage:
```shell
lfr start [flags]
```

## Flags:
- `-h`, `--help`
  - help for `lfr start`

## Global Flags:
- `--no-color`
  - disable colors for output messages