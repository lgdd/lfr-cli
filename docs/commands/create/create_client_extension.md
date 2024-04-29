---
layout: default
title: lfr create client-extension
grand_parent: Commands
parent: lfr create
nav_order: 2
permalink: /commands/create/client-extension
---

# lfr create client-extension

It allows you to create a Client Extension from an official sample. Since Client Extensions are under intensive development, they're not available as templates. So you will need to update the configuration files of any Client Extension you create.


The configuration directory contains the Client Extensions samples. The samples are cloned from a [custom repository](https://github.com/lgdd/liferay-client-extensions-samples){:target="_blank"} that mirrors the [official samples contained in the monorepo](https://github.com/liferay/liferay-portal/tree/master/workspaces/liferay-sample-workspace/client-extensions){:target="_blank"} and rename some components to be more intelligible.

## Usage:
```shell
lfr create workspace NAME [flags]
# or
lfr c ws NAME [flags]
```

Running this command  with no argument triggers the interactive mode where you can choose the template and enter a name from the terminal (and other options depending on the template). To make the interactive mode accessible, edit the `/.lfr/config.toml`:

```toml
[output]
accessible = true
```

## Examples:
```shell
lfr c cx batch my-batch
lfr c cx custom-element-angular my-angular
lfr c cx custom-element-react-dom my-react
lfr c cx etc-node my-node-microservice
lfr c cx etc-spring-boot my-spring-boot-microservice
lfr c cx site-initializer my-site-initializer
lfr c cx commerce-payment-integration my-payment-connector
```

## Flags:
- `-h`, `--help`
  - help for `lfr create`

## Global Flags:
- `--no-color`
  - disable colors for output messages
- `-p`, `--package string`
  - base package name (default "org.acme")