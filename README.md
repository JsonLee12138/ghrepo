# ghrepo

English | [中文](./README.zh.md)

CLI for GitHub repository contents. Browse, inspect, download, create, update, and delete files in any GitHub repository without cloning.

## Table of Contents

- [Installation](#installation)
  - [Homebrew (macOS)](#homebrew-macos)
  - [From Source](#from-source)
  - [From GitHub Releases](#from-github-releases)
  - [Agent Skill](#agent-skill)
- [Upgrade](#upgrade)
  - [Homebrew](#homebrew)
  - [From Source](#from-source-1)
  - [From GitHub Releases](#from-github-releases-1)
  - [Agent Skill](#agent-skill-1)
- [Authentication](#authentication)
  - [Creating a GitHub Personal Access Token](#creating-a-github-personal-access-token)
  - [Using Your Token](#using-your-token)
  - [Best Practices](#best-practices)
  - [Verify Authentication](#verify-authentication)
- [Usage](#usage)
  - [Initialize project](#initialize-project)
  - [List directory contents](#list-directory-contents)
  - [Show file or directory metadata](#show-file-or-directory-metadata)
  - [Output file content](#output-file-content)
  - [Download files or directories](#download-files-or-directories)
  - [Create or update a file](#create-or-update-a-file)
  - [Delete a file](#delete-a-file)
- [Global Flags](#global-flags)
- [License](#license)

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

### Creating a GitHub Personal Access Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Click "Generate new token"
3. Give your token a descriptive name (e.g., "ghrepo CLI")
4. Select the following scopes:
   - `public_repo` - For reading public repositories
   - `repo` - For reading/writing private repositories (if needed)
   - `gist` - Optional, for gist access
5. Click "Generate token" and copy the token immediately
6. **Keep your token secure** - treat it like a password

### Using Your Token

```bash
# Option 1: Set as environment variable (recommended)
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
ghrepo cat owner/repo README.md

# Option 2: Pass directly via flag (for single commands)
ghrepo cat owner/repo README.md --token ghp_xxxxxxxxxxxx

# Option 3: Use GH_TOKEN environment variable
export GH_TOKEN="ghp_xxxxxxxxxxxx"
ghrepo cat owner/repo README.md
```

### Best Practices

- **Use environment variables** - Set `GITHUB_TOKEN` or `GH_TOKEN` in your shell profile for convenience
- **Keep tokens secret** - Never commit tokens to version control
- **Use `.bashrc` or `.zshrc`** - Add to your shell configuration for persistence across sessions:
  ```bash
  # Add to ~/.bashrc or ~/.zshrc
  export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
  ```
- **Rotate tokens periodically** - Update your tokens regularly for security
- **Use token-specific names** - Create separate tokens for different purposes or machines

### Verify Authentication

Verify your token is valid and has the required permissions:

```bash
ghrepo auth check
```

This command will confirm that your token is working and show you what permissions are associated with it.

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
