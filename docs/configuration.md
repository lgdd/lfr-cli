---
layout: page
title: Configuration
nav_order: 3
permalink: /configuration
---

# Configuration
{: .d-inline-block }
v3.0.0+
{: .label .label-green .mb-5 }

`lfr` creates a configuration directory under your home directory: `~/.lfr`. This directory contains a `config.toml` file that allows to configure defaults for some flags:

```toml
[docker]
jdk = 8 # 11, 17 or 21
multistage = false

[logs]
follow = false

[module]
package = 'org.acme' # package name used for your workspace and as default base for your modules

[output]
accessible = false # make command line prompts accessible
no_color = false # remove colors for all outputs

[workspace]
build = 'gradle' # or 'maven'
edition = 'portal' # or 'dxp'
init = false
version = '7.4' # or 7.3, 7.2, 7.1, 7.0
```
> See also [`lfr config`](/cmd/config)

The configuration directory will also contains the Client Extensions samples. The samples are cloned from a [custom repository](https://github.com/lgdd/liferay-client-extensions-samples){:target="_blank"} that mirrors the [official samples contained in the monorepo](https://github.com/liferay/liferay-portal/tree/master/workspaces/liferay-sample-workspace/client-extensions){:target="_blank"} and rename some components to be more intelligible.

The Client Extensions samples are cloned when running the `lfr create client-extension` (or the short version `lfr c cx`) for the first time, and the next time it will refresh it with a `git pull` in the background.