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

Checkout the [release page](https://github.com/lgdd/liferay-cli/releases) to download the binary for your distribution.

For macOS, you can install it using [Homebrew](https://brew.sh):
```shell
brew tap lgdd/homebrew-tap
brew install liferay-cli
```

### Examples:

Get a completion script for your shell:

```shell
lfr completion bash
```

> bash, zsh, fish and powershell are supported.

Run a diagnosis:
```shell
lfr diagnose
# or shorter
lfr diag
```

> The output of this command will list your installations of Java, Blade and Docker. It will also display how much space are being used by cached bundles and docker images.

Create a Liferay workspace:

```shell
lfr create workspace liferay-workspace
# or shorter
lfr create ws liferay-workspace
# or even shoter
lfr c ws liferay-workspace
```

Create a client extension:
```shell
lfr create client-extension
# or shorter
lfr create cx
# or event shorter
lfr c cx
```
> Since client extensions are available as samples in liferay-portal repo, it will checkout the subdirectory containing them under `$HOME/.lfr/liferay-portal`.

Run a Gradle or Maven task:

```shell
lfr exec [tasks...]

# Gradle example
lfr exec build

# Maven example
lfr exec clean install
```

Shortcuts for common tasks:
```shell
# Build
lfr build
#or shorter
lfr b

# Deploy
lfr deploy
# or shorter
lfr d
```

Start Liferay and follow the logs:

```shell
lfr start
lfr logs -f
```

## Benchmarks

Using [Hyperfine](https://github.com/sharkdp/hyperfine).

### Create Workspace

| Command                               |      Mean [s] | Min [s] | Max [s] |     Relative |
| :------------------------------------ | ------------: | ------: | ------: | -----------: |
| `blade init -v 7.4 liferay-workspace` | 1.837 ± 0.300 |   1.665 |   2.668 | 19.94 ± 5.69 |
| `lfr create ws liferay-workspace`     | 0.092 ± 0.022 |   0.076 |   0.178 |         1.00 |

### Create MVC Portlet

| Command                                      |      Mean [s] | Min [s] | Max [s] |       Relative |
| :------------------------------------------- | ------------: | ------: | ------: | -------------: |
| `blade create -t mvc-portlet my-mvc-portlet` | 1.608 ± 0.021 |   1.570 |   1.647 | 59.70 ± 112.37 |
| `lfr create mvc my-mvc-portlet`              | 0.027 ± 0.051 |   0.015 |   0.345 |           1.00 |

### Create Service Builder

| Command                                              |      Mean [s] | Min [s] | Max [s] |       Relative |
| :--------------------------------------------------- | ------------: | ------: | ------: | -------------: |
| `blade create -t service-builder my-service-builder` | 1.628 ± 0.057 |   1.573 |   1.772 | 82.00 ± 134.01 |
| `lfr create sb my-service-builder`                   | 0.020 ± 0.032 |   0.014 |   0.332 |           1.00 |

## License

[MIT](LICENSE)
