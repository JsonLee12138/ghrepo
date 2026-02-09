# ghrepo

CLI for GitHub repository contents. Browse, inspect, download, create, update, and delete files in any GitHub repository without cloning.

## Installation

### Homebrew (macOS)

```bash
brew tap JsonLee12138/ghrepo
brew install --cask ghrepo
```

### From Source

Requires Go 1.22+.

```bash
go install githubRAGCli/cmd/ghrepo@latest
```

### From GitHub Releases

Download the binary for your platform from [Releases](https://github.com/JsonLee12138/ghrepo/releases), extract it, and add it to your `PATH`.

### Agent Skill

Install ghrepo as an [Agent Skill](https://agentskills.io/) for Claude Code, Cursor, Codex and other AI coding agents:

```bash
npx skills add JsonLee12138/ghrepo
```

## Upgrade

### Homebrew

```bash
brew update
brew upgrade --cask ghrepo
```

### From Source

```bash
go install githubRAGCli/cmd/ghrepo@latest
```

### From GitHub Releases

Download the latest version from [Releases](https://github.com/JsonLee12138/ghrepo/releases) and replace the existing binary.

### Agent Skill

```bash
npx skills add JsonLee12138/ghrepo
```

## Authentication

ghrepo requires a GitHub personal access token. You can provide it in three ways (in order of precedence):

1. `--token` flag
2. `GITHUB_TOKEN` environment variable
3. `GH_TOKEN` environment variable

```bash
# Set token via environment variable
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"

# Or pass it directly
ghrepo cat owner/repo README.md --token ghp_xxxxxxxxxxxx
```

Verify your token is valid:

```bash
ghrepo auth check
```

## Usage

### Initialize project

Generate an `AGENTS.md` file in the current directory, configured for a specific repository:

```bash
ghrepo init owner/repo
```

### List directory contents

```bash
ghrepo ls owner/repo src/
ghrepo ls owner/repo src/ --ref develop
ghrepo ls owner/repo src/ --recursive
```

### Show file or directory metadata

```bash
ghrepo stat owner/repo path/to/file
ghrepo stat owner/repo path/to/file --ref main
```

### Output file content

```bash
ghrepo cat owner/repo README.md
ghrepo cat owner/repo src/main.go --ref v1.0.0
```

### Download files or directories

```bash
ghrepo get owner/repo src/ --out ./local-src
ghrepo get owner/repo README.md --out ./README.md --overwrite
```

### Create or update a file

```bash
# Upload a local file (will prompt for confirmation)
ghrepo put owner/repo path/to/file.txt -m "add file" --file ./local.txt

# Upload from stdin (requires --yes)
echo "hello" | ghrepo put owner/repo file.txt -m "create file" --stdin --yes

# Specify a branch
ghrepo put owner/repo config.yml -m "update config" --file ./config.yml -b develop
```

### Delete a file

```bash
# Delete a file (will prompt for confirmation)
ghrepo rm owner/repo path/to/old-file.txt -m "remove old file"

# Skip confirmation
ghrepo rm owner/repo temp.txt -m "cleanup" --yes
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--token` | GitHub personal access token |
| `--api-base` | GitHub API base URL |
| `--timeout` | HTTP request timeout |
| `--json` | Output in JSON format |
| `--verbose` | Enable verbose output |

## License

MIT
