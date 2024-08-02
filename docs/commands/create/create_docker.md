---
layout: default
title: docker
grand_parent: Commands
parent: create
nav_order: 3
permalink: /cmd/create/docker
---

# lfr create docker

## Usage:
```shell
lfr create docker [flags]
# or
lfr c docker [flags]
```

To change the default Java version and type of build (multi-stage or not), you can edit the `/.lfr/config.toml`:
```toml
[docker]
jdk = 8 # 11, 17 or 21
multistage = true
```

## Flags:
- `-h`, `--help`
  - help for `lfr create docker`
- `-j`, `--java` 
  - Java version to use in the Dockerfile (8, 11, 17 or 21)
- `-m`, `--multi-stage`
  - Creates a multi-stage Dockerfile (build & run stages)

## Global Flags:
- `--no-color`
  - disable colors for output messages