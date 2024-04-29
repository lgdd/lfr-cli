---
layout: default
title: spring-mvc
grand_parent: Commands
parent: create
nav_order: 7
permalink: /cmd/create/spring
---

# lfr create spring-mvc-portlet

It allows you to create a Spring MVC Portlet using [PortletMVC4Spring](https://github.com/liferay/portletmvc4spring).

## Usage:
```shell
lfr create spring-mvc-portlet NAME [flags]
# or
lfr c spring NAME [flags]
```

## Flags:
- `-h`, `--help`
  - help for `lfr create spring-mvc-portlet`
- `-t`, `--template`
  - template engine to use (`thymeleaf` or `jsp`) (default `thymeleaf`)

## Global Flags:
- `--no-color`
  - disable colors for output messages
- `-p`, `--package`
  - base package name (default `org.acme`)