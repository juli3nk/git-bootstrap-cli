# Agent Guide — git-bootstrap-cli

This file contains the conventions and commands agents should know when working on this repository.

## Project overview

A small CLI tool to bootstrap Git repositories with standardized project templates. It keeps boilerplate files in a dedicated templates repository outside of target projects.

## Technology stack

- **Language:** Go `1.26.1`
- **CLI framework:** [spf13/cobra](https://github.com/spf13/cobra)
- **CI / local pipelines:** Dagger `v0.21.9`
- **Release build:** GoReleaser (builds binaries and publishes the GitHub release)
- **Release versioning:** semantic-release (creates the tag and the `CHANGELOG.md`)
- **Dev container:** `mcr.microsoft.com/devcontainers/base:ubuntu-26.04`

## Repository structure

```text
.
cmd/git-bootstrap/        # CLI entrypoint
internal/
  bootstrap/              # Template application, filesystem helpers, reports
  cli/                    # Cobra command trees
  config/                 # Configuration loading
  git/                    # Git helpers
.devcontainer/            # Dev container configuration
  devcontainer.json       # Dev container definition
  devcontainer-lock.json  # Pinned feature versions
dagger.json               # Dagger engine + toolchain pins
.github/workflows/        # GitHub Actions entry points
```

## Development commands

Build the project:

```sh
go build ./...
```

Run tests:

```sh
go test ./...
```

Build a snapshot release with GoReleaser via Dagger:

```sh
dagger call goreleaser build --snapshot directory --path=./dist export --path=./dist
```

The entrypoint binary is built from `cmd/git-bootstrap`.

## Lint / format

We use `golangci-lint` v2 configuration (`.golangci.yml`). Enabled linters:

- `gosec`
- `revive` (with `exported` rule disabled)
- `misspell`
- `staticcheck`
- `errcheck`
- `govet`
- `ineffassign`

The Go feature in the devcontainer provides the Go extension; formatting follows `.editorconfig`.

## CI

CI is triggered on pushes and pull requests targeting `develop`.

It delegates to shared workflows in `juli3nk/ci-shared@main`:

- `changes.yaml` — detects changed files
- `git-checks.yaml` — conventional commits linting and secret scanning
- `lint-generic.yaml` — Markdown/JSON linting
- `go-lint.yaml` — Go lint and formatting
- `go-deps.yaml` — Go module, vulnerability and license checks
- `go-build.yaml` — snapshot build via GoReleaser/Dagger

`vars.GO_VERSION` must be set at the repository level and is passed to Go-specific jobs.

> Note: the shared `go-build.yaml` currently does **not** run tests. Adding unit tests is planned as a follow-up.

## Release flow

Releases run on pushes to `main` or on manual trigger via `workflow_dispatch`.

The release pipeline delegates to shared workflows in `juli3nk/ci-shared@main`:

1. `go-lint.yaml` — lints and formats Go code.
2. `go-deps.yaml` — verifies Go modules, scans vulnerabilities, and checks licenses.
3. `go-build.yaml` — builds a snapshot with GoReleaser to ensure `main` compiles.
4. `release.yaml` — `semantic-release` analyzes commits, bumps the version, creates a Git tag, and commits `CHANGELOG.md`.
5. `go-publish.yaml` — if a new tag was created, `goreleaser release` builds the binaries and publishes the GitHub release.

GoReleaser has `changelog.disable: true` because the changelog is managed by semantic-release.

## Environment variables

- `PROJECT_BOOTSTRAP_HOME` — overrides the default templates directory (`~/.local/share/git-bootstrap`)
- `DAGGER_NO_NAG` — set to `1` in CI and devcontainer to disable Dagger update prompts

## Conventions

- Conventional Commits are required (enforced by `commitlint`)
- Branch model: `develop` for continuous integration, `main` for releases
- The `go.mod` version must stay aligned with `vars.GO_VERSION` and the devcontainer Go feature
