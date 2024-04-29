---
layout: default
title: lfr deploy
parent: Commands
nav_order: 4
permalink: /commands/deploy
---

# lfr deploy

It allows you to deploy the artifacts of your Liferay Workspace to your bundle. It works as a shortcut for the corresponding Gradle or Maven goal, i.e. `./gradlew deploy` or `./mvnw bundle-support:deploy`.

If you go in a subdirectory, it will only deploy the artifacts the subdirectory contains.

## Usage:
```shell
lfr deploy [flags]
# or
lfr d
```

## Flags:
- `-h`, `--help`
  - help for `lfr init`

## Global Flags:
- `--no-color`
  - disable colors for output messages