---
layout: page
title: Getting Started
nav_order: 2
permalink: /getting-started
---

# Installation

## Linux / macOS

You can install it using [Homebrew](https://brew.sh){:target="_blank"}:

```shell
brew tap lgdd/homebrew-tap
brew install --cask lfr-cli
```
> If you previously installed it via `brew install lfr-cli` (formula) or `brew install liferay-cli`, run `brew uninstall lfr-cli` then `brew install --cask lfr-cli` to migrate to the new cask.

## Windows

For Windows, go to the [release page](https://github.com/lgdd/lfr-cli/releases){:target="_blank"} and download the zip file corresponding to your architecture. Then extract `lfr.exe` from the archive and move to the folder of your choice. **Make sure that the chosen folder is included in the `%PATH%` environment variable.**

{: .note }
[Winget](https://learn.microsoft.com/en-us/windows/package-manager/winget/) support is pending approval. Once available, you will be able to install via `winget install lgdd.lfr-cli`.

# Completions

```shell
lfr completion [bash|zsh|fish|powershell]
```

To load completions for [Bash](#bash), [Zsh](#zsh), [fish](#fish) and [PowerShell](#powershell):

## Bash:

```shell
source <(lfr completion bash)
# To load completions for each session, execute once:
# Linux:
lfr completion bash > /etc/bash_completion.d/lfr
# macOS:
lfr completion bash > /usr/local/etc/bash_completion.d/lfr
```

## Zsh:

```shell
# If shell completion is not already enabled in your environment,
# you will need to enable it.  You can execute the following once:
echo "autoload -U compinit; compinit" >> ~/.zshrc
# To load completions for each session, execute once:
lfr completion zsh > "${fpath[1]}/_lfr"
# You will need to start a new shell for this setup to take effect.
```

## fish:

```shell
lfr completion fish | source
# To load completions for each session, execute once:
lfr completion fish > ~/.config/fish/completions/lfr.fish
```

## PowerShell:

```powershell
lfr completion powershell | Out-String | Invoke-Expression
# To load completions for every new session, run:
lfr completion powershell > lfr.ps1
# and source this file from your PowerShell profile.
```
