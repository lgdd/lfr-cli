---
layout: default
title: lfr logs
parent: Commands
nav_order: 7
permalink: /commands/logs
---

# lfr start

It allows you to display the logs from the running Liferay bundle (i.e. `catalina.out`). If you want to always follow the logs without the need to add `-f` to the command, you can change edit `~/.lfr/config.toml`:
```toml
[logs]
follow = true # default to false
```

## Usage:
```shell
lfr logs [flags]
```

## Flags:
- `-f`, `--follow`
  - tail the logs from `catalina.out`
- `-h`, `--help`
  - help for `lfr start`

## Global Flags:
- `--no-color`
  - disable colors for output messages