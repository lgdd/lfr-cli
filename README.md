[![new](https://img.shields.io/badge/NEW-Supports%20Client%20Extension-blueviolet)](https://learn.liferay.com/web/guest/w/dxp/building-applications/client-extensions#types-of-client-extensions)

[![version](https://img.shields.io/github/v/tag/lgdd/lfr-cli)](https://github.com/lgdd/lfr-cli/releases)
[![status](https://img.shields.io/github/actions/workflow/status/lgdd/lfr-cli/test.yml)](https://github.com/lgdd/lfr-cli/actions/workflows/test.yml)
[![report](https://goreportcard.com/badge/github.com/lgdd/lfr-cli)](https://goreportcard.com/report/github.com/lgdd/lfr-cli)
[![license](https://img.shields.io/github/license/lgdd/lfr-cli)](https://github.com/lgdd/lfr-cli/blob/main/LICENSE)
![GitHub last commit](https://img.shields.io/github/last-commit/lgdd/lfr-cli?color=teal&label=latest%20commit)

# LFR

`lfr` is an unofficial CLI tool written in Go that helps you create & manage Liferay projects.

![preview](https://github.com/lgdd/doc-assets/blob/main/liferay-cli/liferay-cli-preview.gif?raw=true)

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

## Usage

Checkout the [documentation](https://lgdd.github.io/lfr-cli/) (work in progress).

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
