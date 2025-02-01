---
layout: default
title: trial
parent: Commands
nav_order: 11
permalink: /cmd/trial
---

# lfr trial
{: .d-inline-block }
NEW (3.2.0)
{: .label .label-green .mb-5 }


It will fetch a DXP trial key from [this repo](https://github.com/lgdd/liferay-product-info/tree/main/dxp-trial) which automatically extract a DXP trial key from the latest Docker image of Liferay DXP.

The DXP trial key is then saved in the current directory of the command as a `trial.xml` file. If this file already exists, it won't override it.

## Usage:
```shell
lfr trial
lfr t
```

## Flags:
- `-h`, `--help`
  - help for `lfr trial`
## Global Flags:
- `--no-color`
  - disable colors for output messages