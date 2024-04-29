---
layout: default
title: lfr create
parent: Commands
nav_order: 2
has_children: true
permalink: /commands/create
---

# lfr create

It allows you to create many types of components: Liferay Workspace, Client Extensions, Docker (Dockerfile & docker-compose.yml) and OSGi components (API, Gogo Shell Command, MVC Portlet, Spring Portlet, REST Builder & Service Builder).

## Usage:
```shell
lfr create TYPE NAME [flags]
# or
lfr c TYPE NAME [flags]
```

Running this command  with no argument triggers the interactive mode where you can choose the template and enter a name from the terminal (and other options depending on the template). To make the interactive mode accessible, edit `/.lfr/config.toml`:

```toml
[output]
accessible = true
```

## Available Commands:
- `lfr create workspace`
- `lfr create client-extension`
- `lfr create docker`
- `lfr create api`
- `lfr create command`
- `lfr create mvc-portlet`
- `lfr create spring-mvc-portlet`
- `lfr create rest-builder`
- `lfr create service-builder` 

## Flags:
- `-h`, `--help`
  - help for `lfr create`
- `-p`, `--package string`
  - base package name (default "org.acme")

## Global Flags:
- `--no-color`
  - disable colors for output messages