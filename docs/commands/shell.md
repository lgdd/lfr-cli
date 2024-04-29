---
layout: default
title: shell
parent: Commands
nav_order: 10
permalink: /cmd/shell
---

# lfr shell

{: .warning }
**[Known issue](https://github.com/lgdd/lfr-cli/issues/3)**: The `shell` command is very slow to connect.

It allows you to connect and get the Gogo Shell. The shell returned by this command to communicate to Liferay Gogo Shell allows to use the keyword `exit` safely as it will disconnect instead of stopping the OSGi container. `Ctrl+C` also disconnect from the shell.

## Usage:
```shell
lfr shell [flags]
lfr sh [flags]
```

## Flags:
- `-h`, `--help`
  - help for `lfr stop`

- `--host`
  - default is "localhost"
- `-p`, `--port`
  -  default is 11311
## Global Flags:
- `--no-color`
  - disable colors for output messages