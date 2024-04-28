[![new](https://img.shields.io/badge/NEW-Supports%20Client%20Extension-blueviolet)](https://learn.liferay.com/web/guest/w/dxp/building-applications/client-extensions#types-of-client-extensions)

[![version](https://img.shields.io/github/v/tag/lgdd/lfr-cli)](https://github.com/lgdd/lfr-cli/releases)
[![status](https://img.shields.io/github/actions/workflow/status/lgdd/lfr-cli/test.yml)](https://github.com/lgdd/lfr-cli/actions/workflows/test.yml)
[![report](https://goreportcard.com/badge/github.com/lgdd/lfr-cli)](https://goreportcard.com/report/github.com/lgdd/lfr-cli)
[![license](https://img.shields.io/github/license/lgdd/lfr-cli)](https://github.com/lgdd/lfr-cli/blob/main/LICENSE)
![GitHub last commit](https://img.shields.io/github/last-commit/lgdd/lfr-cli?color=teal&label=latest%20commit)

# LFR

`lfr` is an unofficial CLI tool written in Go that helps you create & manage Liferay projects.

![preview](https://github.com/lgdd/doc-assets/blob/main/liferay-cli/liferay-cli-preview.gif?raw=true)

**Table of contents**:

- [LFR](#lfr)
  - [Motivation](#motivation)
  - [Installation](#installation)
    - [Completions](#completions)
      - [Bash:](#bash)
      - [Zsh:](#zsh)
      - [fish:](#fish)
      - [PowerShell:](#powershell)
  - [Usage](#usage)
    - [diagnose](#diagnose)
    - [create](#create)
      - [workspace](#workspace)
      - [client-extension](#client-extension)
      - [api](#api)
      - [service-builder](#service-builder)
      - [rest-builder](#rest-builder)
      - [mvc-portlet](#mvc-portlet)
      - [spring-mvc-portlet](#spring-mvc-portlet)
      - [docker](#docker)
      - [exec](#exec)
      - [build](#build)
      - [deploy](#deploy)
      - [init](#init)
      - [start](#start)
      - [stop](#stop)
      - [logs](#logs)
      - [shell](#shell)
  - [Benchmarks](#benchmarks)
    - [Create Workspace (Gradle)](#create-workspace-gradle)
    - [Create MVC Portlet](#create-mvc-portlet)
    - [Create Service Builder](#create-service-builder)
  - [License](#license)

## Motivation

I needed a subject to play with Go. Writing a CLI tool is fun - especially with [Cobra](https://github.com/spf13/cobra) - and I wanted to explore how to distribute it using GitHub Actions (and [GoReleaser](https://github.com/goreleaser/goreleaser)).

Also, I get sometimes frustrated by [Blade](https://github.com/liferay/liferay-blade-cli) and wanted to focus on providing:

- Better performances (cf. [benchmarks](benchmarks))
- Better support for Maven
- Shorter commands
- More consistent commands names and ordering
- Details after any command execution
- Shell completion

I'm not the only one motivated to help Liferay developpers with new dev tools. If you're looking for something to help you with Client Extensions development, definitely checkout this tool: https://github.com/bnheise/ce-cli

## Installation

For macOS or Linux, you can install it using [Homebrew](https://brew.sh):

```shell
brew tap lgdd/homebrew-tap
brew install lfr-cli
```
> if you already installed it using `brew install liferay-cli`, make sure to run `brew update` to be able to update to the latest version matching the new name. 

For Windows, go to the [release page](https://github.com/lgdd/lfr-cli/releases) and download the zip file corresponding to your architecture. Then extract `lfr.exe` from the archive and move to the folder of your choice. **Make sure that the chosen folder is included in the `%PATH%` environment variable.**

### Completions

Usage:

```shell
lfr completion [bash|zsh|fish|powershell]
```

To load completions for [Bash](#bash), [Zsh](#zsh), [fish](#fish) and [PowerShell](#powershell):

#### Bash:

```shell
source <(lfr completion bash)
# To load completions for each session, execute once:
# Linux:
lfr completion bash > /etc/bash_completion.d/lfr
# macOS:
lfr completion bash > /usr/local/etc/bash_completion.d/lfr
```

#### Zsh:

```shell
# If shell completion is not already enabled in your environment,
# you will need to enable it.  You can execute the following once:

echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
lfr completion zsh > "${fpath[1]}/_lfr"

# You will need to start a new shell for this setup to take effect.
```

#### fish:

```shell
lfr completion fish | source

# To load completions for each session, execute once:
lfr completion fish > ~/.config/fish/completions/lfr.fish
```

#### PowerShell:

```powershell
PS> lfr completion powershell | Out-String | Invoke-Expression

# To load completions for every new session, run:
PS> lfr completion powershell > lfr.ps1
# and source this file from your PowerShell profile.
```

## Usage

### diagnose

Run a diagnosis to verify your environment for Liferay development:

```shell
lfr diag
```

```shell
Aliases:
  diagnose, diag
```

The output of this command will list your installations of Java, Blade and Docker. It will also display how much space are being used by cached bundles and docker images.

Example:
![diag-example](https://github.com/lgdd/doc-assets/blob/main/lfr-cli/lfr-diag-example.png?raw=true)

### create

#### workspace

```shell
lfr create workspace liferay-workspace
```

```shell
Aliases:
  workspace, ws

Flags:
  -b, --build string     build tool (gradle or maven) (default "gradle")
  -h, --help             help for workspace
  -i, --init             executes Liferay bundle initialization (i.e. download & unzip in the workspace)
  -v, --version string   Liferay major version (7.x) (default "7.4")
```

#### client-extension

Client extensions extend Liferay (7.4 U45+/GA45+) without using OSGi modules.

```shell
  lfr create client-extension
```

```shell
Aliases:
  client-extension, cx

Flags:
  -h, --help   help for client-extension
```

Since client extensions are available as [samples in liferay-portal repo](https://github.com/liferay/liferay-portal/tree/master/workspaces/liferay-sample-workspace/client-extensions), it will checkout the subdirectory containing them under `$HOME/.lfr/liferay-portal`.

#### api

Creates a minimal OSGi module, with an `Export-Package` directive:

```shell
lfr create api my-api
```

```shell
Flags:
  -h, --help   help for api

Global Flags:
      --no-color         disable colors for output messages
  -p, --package string   base package name (default "org.acme")
```

#### service-builder

```shell
lfr create service-builder my-service-builder
```

```shell
Aliases:
  service-builder, sb

Flags:
  -h, --help   help for service-builder

Global Flags:
      --no-color         disable colors for output messages
  -p, --package string   base package name (default "org.acme")
```

#### rest-builder

```shell
lfr create rest-builder my-rest-service
```

```shell
Aliases:
  rest-builder, rb

Flags:
  -g, --generate   executes code generation
  -h, --help       help for rest-builder

Global Flags:
      --no-color         disable colors for output messages
  -p, --package string   base package name (default "org.acme")
```

#### mvc-portlet

```shell
lfr create mvc-portlet my-portlet
```

```shell
Aliases:
  mvc-portlet, mvc

Flags:
  -h, --help   help for mvc-portlet

Global Flags:
      --no-color         disable colors for output messages
  -p, --package string   base package name (default "org.acme")
```

#### spring-mvc-portlet

```shell
lfr create spring-mvc-portlet my-spring-mvc
```

```shell
Aliases:
  spring-mvc-portlet, spring

Flags:
  -h, --help              help for spring-mvc-portlet
  -t, --template string   template engine (thymeleaf or jsp) (default "thymeleaf")

Global Flags:
      --no-color         disable colors for output messages
  -p, --package string   base package name (default "org.acme")
```

#### docker

Creates Dockerfile (multi-stag or not) and docker-compose samples:

```shell
lfr create docker
```

```shell
Flags:
  -h, --help          help for docker
  -j, --java int      Java version (8 or 11) (default 11)
  -m, --multi-stage   use multi-stage build

Global Flags:
      --no-color         disable colors for output messages
  -p, --package string   base package name (default "org.acme")
```

#### exec

Execute Gradle or Maven task(s):

```shell
# Gradle
lfr exec build

# Maven
lfr exec clean install
```

```shell
Aliases:
  exec, x

Flags:
  -h, --help   help for exec

Global Flags:
      --no-color   disable colors for output messages
```

#### build

Shortcut to build your Liferay bundle:

```shell
lfr build
```

```shell
Aliases:
  build, b

Flags:
  -h, --help   help for build

Global Flags:
      --no-color   disable colors for output messages
```

#### deploy

Shortcut to deploy your modules using Gradle or Maven:

```shell
lfr deploy
```

```shell
Aliases:
  deploy, d

Flags:
  -h, --help   help for deploy

Global Flags:
      --no-color   disable colors for output messages
```

#### init

Shortcut to initialize your Liferay bundle:

```shell
lfr init
```

```shell
Aliases:
  init, i

Flags:
  -h, --help   help for init

Global Flags:
      --no-color   disable colors for output messages
```

#### start

Start the Liferay Tomcat bundle initialized in your workspace:

```shell
lfr start
```

> It creates a file named `liferay.pid` which is used by `lfr status` display if Tomcat is running or not (or it can be used to kill it if necessary).

```shell
Flags:
  -h, --help   help for start

Global Flags:
      --no-color   disable colors for output messages
```

#### stop

Stop the Liferay Tomcat bundle initialized in your workspace:

```shell
lfr stop
```
Flags:
  -h, --help   help for stop

Global Flags:
      --no-color   disable colors for output messages
```

#### status

Display the status (running or stopped) of Liferay Tomcat from your workspace:

```shell
lfr status
```

```shell
Flags:
  -h, --help   help for start

Global Flags:
      --no-color   disable colors for output messages
```

#### logs

Display and follow the logs from the running Liferay bundle:

```shell
lfr logs -f
```

```shell
Aliases:
  logs, l

Flags:
  -f, --follow   --follow
  -h, --help     help for logs

Global Flags:
      --no-color   disable colors for output messages
```

#### shell

Connect and get Liferay Gogo Shell:

```shell
lfr sh
```

```shell
Aliases:
  shell, sh

Flags:
  -h, --help          help for shell
      --host string   --host localhost (default "localhost")
  -p, --port int      --port 11311 (default 11311)

Global Flags:
      --no-color   disable colors for output message
```

The keyword `exit` can be safely used with this shell as it will disconnect instead of stopping the OSGi container. `Ctrl+C` also disconnect from the shell.

**Known issue**: [The "shell" command is very slow to connect](https://github.com/lgdd/lfr-cli/issues/3)

## Benchmarks

Using [Hyperfine](https://github.com/sharkdp/hyperfine) with 5 warmup runs and the following setup:
- `blade version 6.0.0.202404102137`
- `lfr version v3.0.0`

### Create Workspace (Gradle)

| Command                               |      Mean [s] | Min [s] | Max [s] |     Relative |
| :------------------------------------ | ------------: | ------: | ------: | -----------: |
| `lfr c ws liferayws`                  | 179.5 ± 17.2  | 165.9   | 237.1   | 1.00         |
| `blade init -v 7.4 liferayws`         | 719.6 ± 7.2   | 709.5   | 728.9   | 4.01 ± 0.39  |
> Summary: `lfr c ws liferayws` ran 4.01 ± 0.39 times faster than `blade init -v 7.4 liferayws`

### Create MVC Portlet

| Command                                      |      Mean [s] | Min [s] | Max [s] |       Relative |
| :------------------------------------------- | ------------: | ------: | ------: | -------------: |
| `lfr c mvc my-mvc-portlet`                   | 23.0 ± 4.4    | 20.5    | 47.0    | 1.00           |
| `blade create -t mvc-portlet my-mvc-portlet` | 710.3 ± 6.1   | 703.6   | 720.5   | 30.82 ± 5.94   |
> Summary: `lfr c mvc my-mvc-portlet` ran 30.82 ± 5.94 times faster than `blade create -t mvc-portlet my-mvc-portlet`

### Create Service Builder

| Command                                              |      Mean [s] | Min [s] | Max [s] |       Relative |
| :--------------------------------------------------- | ------------: | ------: | ------: | -------------: |
| `lfr create sb my-service-builder`                   | 19.8 ± 2.8    | 17.6    | 31.2    | 1.00           |
| `blade create -t service-builder my-service-builder` | 723.8 ± 11.0  | 708.9   | 738.8   | 36.65 ± 5.17   |
> Summary: `lfr create sb my-service-builder` ran 36.65 ± 5.17 times faster than `blade create -t service-builder my-service-builder`

## License

[MIT](LICENSE)
