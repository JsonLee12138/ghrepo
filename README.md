# ghrepo

Read-only CLI for GitHub repository contents. Browse, inspect, and download files from any GitHub repository without cloning.

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
