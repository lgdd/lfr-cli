---
layout: default
title: rest-builder
grand_parent: Commands
parent: create
nav_order: 8
permalink: /cmd/create/rest-builder
---

# lfr create rest-builder

It allows you to create a REST Builder module to expose REST APIs & GraphQL endpoints using OpenAPI v3.

## Usage:
```shell
lfr create rest-builder NAME [flags]
# or
lfr c rb NAME [flags]
```

## Flags:

- `-g`, `--generate`
  - executes code generation
- `-h`, `--help`
  - help for `lfr create rest-builder`

## Global Flags:
- `--no-color`
  - disable colors for output messages
- `-p`, `--package`
  - base package name (default `org.acme`)