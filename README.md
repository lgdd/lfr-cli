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
|:---|---:|---:|---:|---:|
| `blade init -v 7.4.1-1 liferay-workspace` | 1.512 ± 0.036 | 1.466 | 1.577 | 115.60 ± 8.57 |
| `lfr create ws liferay-workspace` | 0.013 ± 0.001 | 0.011 | 0.016 | 1.00 |

### Create MVC Portlet

| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `blade create -t mvc-portlet my-mvc-portlet` | 1.628 ± 0.055 | 1.576 | 1.750 | 105.73 ± 18.09 |
| `lfr create mvc my-mvc-portlet` | 0.015 ± 0.003 | 0.012 | 0.027 | 1.00 |

### Create Service Builder

| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `blade create -t service-builder my-service-builder` | 1.736 ± 0.322 | 1.605 | 2.651 | 130.08 ± 38.89 |
| `lfr create sb my-service-builder` | 0.013 ± 0.003 | 0.011 | 0.035 | 1.00 |

## License

[MIT](LICENSE)
