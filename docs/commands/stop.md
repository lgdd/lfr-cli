---
layout: default
title: stop
parent: Commands
nav_order: 8
permalink: /cmd/stop
---

# lfr stop

It allows you to stop the Tomcat bundle started from your Liferay Workspace. To stop the bundle, this command needs the `liferay.pid` file created at the root of the Liferay Workspace by the command `lfr start`.

While waiting for the process to terminate, a spinner displays the PID. Once done, a success or warning message is printed depending on whether the process has fully stopped.

## Usage:
```shell
lfr stop [flags]
```

## Flags:
- `-h`, `--help`
  - help for `lfr stop`

## Global Flags:
- `--no-color`
  - disable colors for output messages