---
layout: default
title: config
parent: Commands
nav_order: 2
permalink: /cmd/config
---

# lfr config
{: .d-inline-block }
v3.1.0+
{: .label .label-info .mb-5 }

Helper command to quickly set new values in your configuration as an alternative to edit the file under `$HOME/.lfr/config.toml` (see [Configuration](/configuration)).

## Usage:
```shell
lfr config -l # see the list of key values stored in the configuration
lfr config key # get the value of a key
lfr config key value # set a new key value with a space
lfr config key=value # or with an equal sign
```

## Flags:
- `-h`, `--help`
  - help for `lfr diagnose`
- `-l`, `--list`
  - list the key/values for your current configuration

## Global Flags:
- `--no-color`
  - disable colors for output messages