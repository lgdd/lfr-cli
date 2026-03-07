# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build and install
go install ./...

# Run all checks (tidy, verify, fmt, vet, test, install)
make

# Lint (requires golangci-lint and gocyclo)
make lint

# Run tests with coverage report
make test

# Run a single test
go test ./pkg/scaffold/... -run TestCreateWorkspace_Gradle_ShouldHaveExpectedFiles -v

# Run tests for a specific package
go test ./pkg/scaffold/... -v
```

## Architecture

The CLI is built with [Cobra](https://github.com/spf13/cobra) and follows a layered structure:

```
cmd/lfr/main.go          # Entry point
internal/cmd/root.go     # Cobra root, registers all subcommands
internal/cmd/<cmd>/      # One package per subcommand (build, create, deploy, exec, etc.)
internal/assets/         # Embedded templates via go:embed (tpl/)
internal/conf/           # Config management via Viper (~/.lfr/config.toml)
internal/prompt/         # Interactive TUI forms (charmbracelet/huh)
pkg/scaffold/            # Core scaffolding logic that generates project files
pkg/metadata/            # Workspace metadata, fetches release info from GitHub
pkg/util/fileutil/       # Template processing (text/template) and workspace detection
pkg/util/helper/         # Client extension sample management (clones from GitHub)
pkg/util/logger/         # Styled logging via charmbracelet/lipgloss
pkg/util/procutil/       # Shell command execution and Java/process detection
```

### Key Flows

**Template system**: Project templates live in `internal/assets/tpl/` and are embedded into the binary via `go:embed`. Files are first copied verbatim by `fileutil.CreateFilesFromAssets`, then rendered through Go's `text/template` with `fileutil.UpdateWithData`.

**Workspace detection**: The current directory's Liferay workspace is located by walking up parent directories until `platform.bndrun` is found (`fileutil.GetLiferayWorkspacePath`).

**Client extensions**: Samples are not bundled in the binary. On first use, `pkg/util/helper` clones `lgdd/liferay-client-extensions-samples` and `lgdd/liferay-client-extensions-extra-samples` into `~/.lfr/`. Subsequent invocations run `git pull` in the background. If GitHub is unreachable, the embedded fallback in `internal/assets/tpl/client_extension/` is used instead.

**Release metadata**: `pkg/metadata` fetches Liferay version info (bundle URLs, Docker images, target platform) from `lgdd/liferay-product-info` on GitHub. Three dispatch paths exist in `NewWorkspaceData`:
- **DXP quarterly** (`YYYY.qN`, e.g. `2024.q1`): loads `~/.lfr/releases.json` from disk on subsequent calls (triggering a background refresh via `refreshReleasesJSON`), or fetches and saves it from GitHub on first use (`fetchAndWriteReleasesJSON`). Filters by `productGroupVersion`, picks the best release (preferred: `recommended` tag → `promoted == "true"` → first entry), then fetches `release.properties` from the CDN. Offline fallback: `getOfflineDXPQuarterlyData` (hardcoded data for 2023.q3–2025.q4).
- **Portal specific GA** (`7.4.3.N-gaN`, e.g. `7.4.3.112-ga112`): fetches `release.properties` directly from the CDN. Offline fallback: `getOfflinePortalGAData` (hardcoded data for ga98–ga132).
- **Legacy 7.x** (`7.0`–`7.4`): fetches the per-version `releases/{edition}_{version}_releases.json`. Offline fallback: `getOfflineWorkspaceData`.

**Interactive vs non-interactive**: Every `create` subcommand works both with direct CLI arguments and with a charmbracelet/huh TUI prompt when no arguments are given. The prompt path re-executes the command by appending synthesized args to `os.Args`.

**Configuration**: Stored in `~/.lfr/config.toml`, managed by Viper. Keys are constants in `internal/conf/conf.go` (e.g. `workspace.edition`, `workspace.version`, `workspace.build`, `module.package`). Defaults are set via `conf.setDefaults()` on first run.

### Conventions

**Scaffold error handling**: All public functions in `pkg/scaffold/` return `error` and never call `logger.Fatal` directly. The cmd layer (`internal/cmd/create/`) is responsible for handling errors and calling `logger.Fatal`. Shared helpers live in `pkg/scaffold/helpers.go`.

**Testing scaffold functions**: Functions that call `fileutil.GetLiferayWorkspacePath()` internally (e.g. `CreateModuleAPI`, `CreateModuleMVC`, `CreateModuleGogoCommand`, `CreateModuleSpring`) require the working directory to be inside a Liferay workspace. Use the `chdirWorkspace(t, path)` test helper (defined in `helpers_test.go`) before calling them. Functions that take the workspace path as a parameter (e.g. `CreateModuleRESTBuilder`, `CreateModuleServiceBuilder`) do not need this.
