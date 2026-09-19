# git-bootstrap-cli

A small CLI tool to bootstrap Git repositories with standardized project templates.

It keeps your boilerplate (licenses, CI configs, base files, stack-specific files, etc.) outside of your projects in a dedicated templates repository. You can install it once and reuse it across many repositories.

## Features

- Install and update a templates repository from any Git remote.
- Create new projects from layered templates: `common/`, `stacks/<stack>/base/`, optional stack kinds, and append-only snippets.
- Upgrade existing projects with the latest template changes.
- Compare a project against its expected templates with `doctor`.
- Create the initial bootstrap commits and tag according to a `templates.yml` manifest.

## Installation

### From GitHub releases

Download the latest binary for your platform from the [releases page](https://github.com/juli3nk/git-bootstrap-cli/releases) and place it somewhere on your `PATH`.

### With `go install`

```sh
go install github.com/juli3nk/git-bootstrap-cli/cmd/git-bootstrap@latest
```

## Configuration

By default, templates are stored in:

```text
~/.local/share/git-bootstrap/templates
```

You can override this location with the `PROJECT_BOOTSTRAP_HOME` environment variable:

```sh
export PROJECT_BOOTSTRAP_HOME=/path/to/your/bootstrap/home
```

## Commands

| Command | Description |
|---------|-------------|
| `git-bootstrap templates install <repo-url>` | Clone a templates repository into the configured templates directory. |
| `git-bootstrap templates update` | Pull the latest changes from the templates repository. |
| `git-bootstrap templates status` | Show the templates directory, version and remote origin. |
| `git-bootstrap new [stack] [kind]` | Apply templates to the current directory. Use `--license` to add a license file. |
| `git-bootstrap upgrade [stack] [kind]` | Re-apply the latest templates to the current directory. |
| `git-bootstrap doctor [stack] [kind]` | Compare current files against the expected templates. |
| `git-bootstrap commit init` | Initialize a Git repository and create the configured bootstrap commits and tag. |
| `git-bootstrap commit app` | Commit staged changes with the configured application commit message. |
| `git-bootstrap version` | Print the current binary version. |

## Templates repository layout

A valid templates repository looks like this:

```text
templates/
├── common/                 # Files applied to every project
├── licenses/
│   ├── MIT.txt
│   └── APACHE-2.0.txt
├── stacks/
│   └── go/
│       ├── base/           # Base files for the "go" stack
│       │   └── append/     # Files appended, not overwritten
│       └── webapp/         # Optional kind under the "go" stack
│           └── append/     # Kind-specific appended files
└── templates.yml           # Manifest for bootstrap commits
```

The `append/` directories are optional. Files inside them are appended to existing files instead of overwriting them.

## Quick example

```sh
# Install a templates repository once
git-bootstrap templates install https://github.com/your-org/bootstrap-templates.git

# Create a new Go web project in the current directory
git-bootstrap new go webapp --license MIT

# Later, check if the project still matches the templates
git-bootstrap doctor go webapp
```

## Development

```sh
go build ./...
go test ./...
```

## License

This project is licensed under the [MIT License](LICENSE).
