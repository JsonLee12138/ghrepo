# ghrepo

[English](./README.md) | 中文

用于 GitHub 仓库内容的 CLI。无需克隆即可浏览、查看、下载、创建、更新和删除任意 GitHub 仓库中的文件。

## 目录

- [安装](#安装)
  - [Homebrew (macOS)](#homebrew-macos)
  - [从源码安装](#从源码安装)
  - [从 GitHub 发行版安装](#从-github-发行版安装)
  - [Agent Skill](#agent-skill)
- [升级](#升级)
  - [Homebrew](#homebrew)
  - [从源码](#从源码)
  - [从 GitHub 发行版](#从-github-发行版)
  - [Agent Skill](#agent-skill-1)
- [认证](#认证)
  - [创建 GitHub 个人访问令牌](#创建-github-个人访问令牌)
  - [使用令牌](#使用令牌)
  - [最佳实践](#最佳实践)
  - [验证认证](#验证认证)
- [使用说明](#使用说明)
  - [初始化项目](#初始化项目)
  - [列出目录内容](#列出目录内容)
  - [查看文件或目录元数据](#查看文件或目录元数据)
  - [输出文件内容](#输出文件内容)
  - [下载文件或目录](#下载文件或目录)
  - [创建或更新文件](#创建或更新文件)
  - [删除文件](#删除文件)
- [全局参数](#全局参数)
- [许可证](#许可证)

## 安装

### Homebrew (macOS)

```bash
brew tap JsonLee12138/ghrepo
brew install --cask ghrepo
```

### 从源码安装

需要 Go 1.22+。

```bash
go install githubRAGCli/cmd/ghrepo@latest
```

### 从 GitHub 发行版安装

从 [Releases](https://github.com/JsonLee12138/ghrepo/releases) 下载对应平台的二进制包，解压后将可执行文件加入 `PATH`。

### Agent Skill

将 ghrepo 安装为 [Agent Skill](https://agentskills.io/)，供 Claude Code、Cursor、Codex 等 AI 编程助手使用：

```bash
npx skills add JsonLee12138/ghrepo
```

## 升级

### Homebrew

```bash
brew update
brew upgrade --cask ghrepo
```

### 从源码

```bash
go install githubRAGCli/cmd/ghrepo@latest
```

### 从 GitHub 发行版

从 [Releases](https://github.com/JsonLee12138/ghrepo/releases) 下载最新版本并替换现有二进制文件。

### Agent Skill

```bash
npx skills add JsonLee12138/ghrepo
```

## 认证

ghrepo 需要 GitHub 个人访问令牌。可通过以下三种方式提供（优先级从高到低）：

1. `--token` 参数
2. 环境变量 `GITHUB_TOKEN`
3. 环境变量 `GH_TOKEN`

### 创建 GitHub 个人访问令牌

1. 打开 [GitHub 设置 > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. 点击「Generate new token」
3. 为令牌起一个易识别的名称（如 "ghrepo CLI"）
4. 勾选以下权限：
   - `public_repo`：读取公开仓库
   - `repo`：读写私有仓库（如需）
   - `gist`：可选，用于 Gist
5. 点击「Generate token」并立即复制令牌
6. **妥善保管令牌**——请像对待密码一样保管

### 使用令牌

```bash
# 方式一：设为环境变量（推荐）
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
ghrepo cat owner/repo README.md

# 方式二：通过参数传入（单次命令）
ghrepo cat owner/repo README.md --token ghp_xxxxxxxxxxxx

# 方式三：使用 GH_TOKEN 环境变量
export GH_TOKEN="ghp_xxxxxxxxxxxx"
ghrepo cat owner/repo README.md
```

### 最佳实践

- **使用环境变量**：在 shell 配置中设置 `GITHUB_TOKEN` 或 `GH_TOKEN` 便于日常使用
- **不要泄露令牌**：切勿将令牌提交到版本库
- **写入 `.bashrc` 或 `.zshrc`**：在 shell 配置中持久化，例如：
  ```bash
  # 添加到 ~/.bashrc 或 ~/.zshrc
  export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
  ```
- **定期轮换令牌**：为安全起见定期更新令牌
- **按用途命名**：为不同用途或机器创建不同的令牌

### 验证认证

确认令牌有效且具备所需权限：

```bash
ghrepo auth check
```

该命令会验证令牌是否可用，并显示当前关联的权限。

## 使用说明

### 初始化项目

在当前目录生成针对指定仓库配置的 `AGENTS.md`：

```bash
ghrepo init owner/repo
```

### 列出目录内容

```bash
ghrepo ls owner/repo src/
ghrepo ls owner/repo src/ --ref develop
ghrepo ls owner/repo src/ --recursive
```

### 查看文件或目录元数据

```bash
ghrepo stat owner/repo path/to/file
ghrepo stat owner/repo path/to/file --ref main
```

### 输出文件内容

```bash
ghrepo cat owner/repo README.md
ghrepo cat owner/repo src/main.go --ref v1.0.0
```

### 下载文件或目录

```bash
ghrepo get owner/repo src/ --out ./local-src
ghrepo get owner/repo README.md --out ./README.md --overwrite
```

### 创建或更新文件

```bash
# 上传本地文件（会提示确认）
ghrepo put owner/repo path/to/file.txt -m "add file" --file ./local.txt

# 从 stdin 上传（需加 --yes）
echo "hello" | ghrepo put owner/repo file.txt -m "create file" --stdin --yes

# 指定分支
ghrepo put owner/repo config.yml -m "update config" --file ./config.yml -b develop
```

### 删除文件

```bash
# 删除文件（会提示确认）
ghrepo rm owner/repo path/to/old-file.txt -m "remove old file"

# 跳过确认
ghrepo rm owner/repo temp.txt -m "cleanup" --yes
```

## 全局参数

| 参数 | 说明 |
|------|------|
| `--token` | GitHub 个人访问令牌 |
| `--api-base` | GitHub API 基础 URL |
| `--timeout` | HTTP 请求超时时间 |
| `--json` | 以 JSON 格式输出 |
| `--verbose` | 启用详细输出 |

## 许可证

MIT
