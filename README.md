[![new](https://img.shields.io/badge/NEW-Supports%20Client%20Extension-blueviolet)](https://learn.liferay.com/web/guest/w/dxp/building-applications/client-extensions#types-of-client-extensions)

[![version](https://img.shields.io/github/v/tag/lgdd/liferay-cli)](https://github.com/lgdd/liferay-cli/releases)
[![status](https://img.shields.io/github/actions/workflow/status/lgdd/liferay-cli/test.yml)](https://github.com/lgdd/liferay-cli/actions/workflows/test.yml)
[![report](https://goreportcard.com/badge/github.com/lgdd/liferay-cli)](https://goreportcard.com/report/github.com/lgdd/liferay-cli)
[![license](https://img.shields.io/github/license/lgdd/liferay-cli)](https://github.com/lgdd/liferay-cli/blob/main/LICENSE)

# Liferay CLI

Liferay CLI - `lfr` - is an unofficial tool written in Go that helps you create & manage Liferay projects.

## Why?

I needed a subject to play with Go. Writing a CLI tool is fun - especially with [Cobra](https://github.com/spf13/cobra) - and I wanted to explore how to distribute it using GitHub Actions (and [GoReleaser](https://github.com/goreleaser/goreleaser)).

Also, I get sometimes frustrated by [Blade](https://github.com/liferay/liferay-blade-cli) and wanted to focus on providing:

- Better performances (cf. [benchmarks](benchmarks))
- Better support for Maven
- Shorter commands
- More consistent commands names and ordering
- Details after any command execution
- Shell completion

## Getting Started

This tool is still in `alpha`, so checkout the [release page](https://github.com/lgdd/liferay-cli/releases) to download the binary for your distribution.

Examples:

Get a completion script for your shell:

```shell
lfr completion bash
```

> bash, zsh, fish and powershell are supported.

Create a Liferay workspace:

```shell
lfr create workspace my-liferay-workspace
# or
lfr create ws my-liferay-workspace
```

Run a Gradle or Maven task:

```shell
# Gradle
lfr exec build

# Maven
lfr exec clean install
```

Start Liferay and follow the logs:

```shell
lfr start
lfr logs -f
```

## Benchmarks

Using [Hyperfine](https://github.com/sharkdp/hyperfine).

### Create Workspace

| Command | Mean [s] | Min [s] | Max [s] | Relative |
| :--- | ---: | ---: | ---: | ---: |
| `blade init -v 7.4 liferay-workspace` |   <span style="color:red">1.837 ± 0.300</span> |   <span style="color:red">1.665</span> |   <span style="color:red">2.668</span> | 19.94 ± 5.69 |
| `lfr create ws liferay-workspace`     | <span style="color:green">0.092 ± 0.022</span> | <span style="color:green">0.076</span> | <span style="color:green">0.178</span> |         1.00 |

### Create MVC Portlet

| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `blade create -t mvc-portlet my-mvc-portlet` | <span style="color:red">1.608 ± 0.021</span> | <span style="color:red">1.570</span> | <span style="color:red">1.647</span> | 59.70 ± 112.37 |
| `lfr create mvc my-mvc-portlet` | <span style="color:green">0.027 ± 0.051</span> | <span style="color:green">0.015</span> | <span style="color:green">0.345</span> | 1.00 |

### Create Service Builder

| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `blade create -t service-builder my-service-builder` | <span style="color:red">1.628 ± 0.057</span> | <span style="color:red">1.573</span> | <span style="color:red">1.772</span> | 82.00 ± 134.01 |
| `lfr create sb my-service-builder` | <span style="color:green">0.020 ± 0.032</span> | <span style="color:green">0.014</span> | <span style="color:green">0.332</span> | 1.00 |

## License

[MIT](LICENSE)
