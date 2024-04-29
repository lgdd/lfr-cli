---
layout: default
title: exec
parent: Commands
nav_order: 9
permalink: /cmd/exec
---

# lfr exec

It allows you to run a command for Gradle or Maven, depending on what you're using in your Liferay Workspace.

## Usage:
```shell
lfr exec TASK... -- [TASK_FLAG]... [flags]
lfr x TASK... -- [TASK_FLAG]... [flags]
```

## Examples:
```shell
lfr exec deploy # for Gradle or Maven
lfr x clean install # for Maven
lfr x package -- --debug # with Maven flags
lfr x build -- --stacktrace # with Gradle flags
```

## Flags:
- `-h`, `--help`
  - help for `lfr stop`

## Global Flags:
- `--no-color`
  - disable colors for output messages