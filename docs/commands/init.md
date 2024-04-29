---
layout: default
title: lfr init
parent: Commands
nav_order: 5
permalink: /commands/init
---

# lfr init

It allows you to initialize your Liferay Workspace. It works as a shortcut for the corresponding Gradle or Maven goal, i.e. `./gradlew initBundle` or `./mvnw bundle-support:init`.

## Usage:
```shell
lfr init [flags]
# or
lfr i
```

## Flags:
- `-h`, `--help`
  - help for `lfr init`

## Global Flags:
- `--no-color`
  - disable colors for output messages