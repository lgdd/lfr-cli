---
layout: page
title: Home
nav_order: 1
---

# lfr-cli
{: .d-inline-block }
v3.0.0
{: .label .mb-5 }

`lfr` is an unofficial CLI tool written in Go that helps you create & manage Liferay projects.
{: .fs-6 .fw-300 }

![preview](https://github.com/lgdd/doc-assets/blob/main/liferay-cli/liferay-cli-preview.gif?raw=true)

# Motivation

I needed a subject to play with Go. Writing a CLI tool is fun - especially with [Cobra](https://github.com/spf13/cobra){:target="_blank"} - and I wanted to explore how to distribute it using GitHub Actions (and [GoReleaser](https://github.com/goreleaser/goreleaser){:target="_blank"}).

Also, I get sometimes frustrated by [Blade](https://github.com/liferay/liferay-blade-cli) and wanted to focus on providing:

- Better performances (cf. [benchmarks](https://github.com/lgdd/lfr-cli?tab=readme-ov-file#benchmarks){:target="_blank"})
- Better support for Maven
- Shorter commands
- More consistent commands names and ordering
- Details after any command execution
- Shell completion

# Related projects

To reach some of the goals of this tool, other projects have been built in order to enhance the developer experience:
- [https://github.com/lgdd/liferay-product-info](https://github.com/lgdd/liferay-product-info)
- [https://github.com/lgdd/liferay-dxp-releases](https://github.com/lgdd/liferay-dxp-releases)
- [https://github.com/lgdd/liferay-client-extensions-samples](https://github.com/lgdd/liferay-client-extensions-samples)