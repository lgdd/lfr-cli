---
layout: default
title: api
grand_parent: Commands
parent: create
nav_order: 4
permalink: /cmd/create/api
---

# lfr create api

It allows you to create a simple OSGi module that exposes an `api` package.

## Usage:
```shell
lfr create api NAME [flags]
# or
lfr c api NAME [flags]
```

## Flags:
- `-h`, `--help`
  - help for `lfr create api`

## Global Flags:
- `--no-color`
  - disable colors for output messages
- `-p`, `--package`
  - base package name (default `org.acme`)