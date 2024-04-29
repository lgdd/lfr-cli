---
layout: default
title: build
parent: Commands
nav_order: 3
permalink: /cmd/build
---

# lfr build

It allows you to build your Liferay Workspace. It works as a shortcut for the corresponding Gradle or Maven goal, i.e. `./gradlew clean install` or `./mvnw clean install`.

If you go in a subdirectory, it will only build the artifacts the subdirectory contains.

## Usage:
```shell
lfr build [flags]
# or
lfr b
```

## Flags:
- `-h`, `--help`
  - help for `lfr init`

## Global Flags:
- `--no-color`
  - disable colors for output messages