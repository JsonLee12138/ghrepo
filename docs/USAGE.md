# ghrepo 使用文档（v0.1 草案）

## 1. 工具目标
`ghrepo` 是一个 CLI 工具，用于通过 GitHub Token 访问和管理仓库中的目录和文件，支持：
- 查看目录结构
- 读取文件内容
- 下载文件或目录
- 查询路径元信息
- 创建或更新文件
- 删除文件

默认面向 GitHub.com，后续可通过 `--api-base` 支持 GitHub Enterprise Server (GHES)。

## 2. 安装方式（规划）

### 2.1 后续推荐（Homebrew）
```bash
brew tap <yourname>/tap
brew install ghrepo
```

### 2.2 本地源码构建（开发期）
```bash
git clone <repo-url>
cd githubRAGCli
go build -o ghrepo ./cmd/ghrepo
```

## 3. 认证配置

### 3.1 环境变量
```bash
export GITHUB_TOKEN=ghp_xxx
```

### 3.2 Token 读取优先级
1. `--token`
2. `GITHUB_TOKEN`
3. `GH_TOKEN`

### 3.3 最小权限建议
- 只读操作：Fine-grained PAT：`Contents: Read`
- 写入/删除操作：Fine-grained PAT：`Contents: Read and Write`

## 4. 命令总览

```bash
ghrepo init <owner/repo>
ghrepo auth check
ghrepo ls <owner/repo> <path> [--ref <ref>] [--recursive] [--json]
ghrepo cat <owner/repo> <path> [--ref <ref>]
ghrepo get <owner/repo> <path> --out <local-path> [--ref <ref>] [--overwrite]
ghrepo stat <owner/repo> <path> [--ref <ref>] [--json]
ghrepo put <owner/repo> <path> -m <msg> (--file <local-path> | --stdin) [-b <branch>] [--yes]
ghrepo rm <owner/repo> <path> -m <msg> [-b <branch>] [--yes]
```

## 5. 详细命令

### 5.0 `init`
在当前目录生成 `AGENTS.md` 文件，写入目标仓库地址及常用命令参考。

示例：
```bash
ghrepo init owner/repo
```

行为说明：
- 在当前工作目录下创建 `AGENTS.md` 文件
- 如果文件已存在，会报错并退出（需手动删除后重新初始化）
- 不需要 Token

### 5.1 `auth check`
检查 Token 是否可用、是否可访问 GitHub API。

示例：
```bash
ghrepo auth check
```

输出（示例）：
```text
auth: ok
user: octocat
rate_limit_remaining: 4978
```

### 5.2 `ls`
列出仓库目录内容。

示例：
```bash
ghrepo ls owner/repo docs --ref main
ghrepo ls owner/repo docs --recursive --json
```

行为说明：
- `path` 可以是 `.`、空目录路径或子目录
- 默认非递归；`--recursive` 时返回完整子树

### 5.3 `cat`
读取文件内容并输出到标准输出。

示例：
```bash
ghrepo cat owner/repo README.md --ref main
```

行为说明：
- 仅用于文件路径，目录路径会报错
- 适合配合重定向：`> local-file`

### 5.4 `get`
下载单文件或目录到本地。

示例：
```bash
ghrepo get owner/repo README.md --out ./downloads/README.md
ghrepo get owner/repo docs --out ./downloads/docs --ref main
```

行为说明：
- 文件下载到 `--out` 指定文件路径
- 目录下载到 `--out` 指定目录路径，保留仓库内相对结构
- 默认不覆盖已有文件，使用 `--overwrite` 强制覆盖

### 5.5 `stat`
查询路径元信息（文件/目录）。

示例：
```bash
ghrepo stat owner/repo README.md
ghrepo stat owner/repo docs --json
```

返回字段（示例）：
- `type` (`file`/`dir`)
- `path`
- `sha`
- `size`（目录为 0 或省略）
- `download_url`（文件时可用）

### 5.6 `put`
创建或更新仓库中的文件。自动检测文件是否存在（创建 vs 更新）。

**执行前会要求二次确认，输入 `y` 确认操作。**

示例：
```bash
# 从本地文件上传
ghrepo put owner/repo path/to/file.txt -m "add file" --file ./local.txt

# 从 stdin 读取内容（需配合 --yes）
echo "hello" | ghrepo put owner/repo file.txt -m "create" --stdin --yes

# 指定目标分支
ghrepo put owner/repo config.yml -m "update config" --file ./config.yml -b develop

# 跳过确认
ghrepo put owner/repo file.txt -m "msg" --file ./f.txt --yes
```

参数说明：
- `-m` / `--message`：提交信息（必填）
- `--file <path>`：从本地文件读取内容
- `--stdin`：从标准输入读取内容（与 `--file` 互斥）
- `-b` / `--branch`：目标分支（可选，默认为仓库默认分支）
- `-y` / `--yes`：跳过确认提示

### 5.7 `rm`
删除仓库中的文件。

**执行前会要求二次确认，输入 `y` 确认操作。**

示例：
```bash
# 删除文件（会提示确认）
ghrepo rm owner/repo old-file.txt -m "remove old file"

# 指定分支并跳过确认
ghrepo rm owner/repo temp.txt -m "cleanup" -b develop --yes
```

参数说明：
- `-m` / `--message`：提交信息（必填）
- `-b` / `--branch`：目标分支（可选）
- `-y` / `--yes`：跳过确认提示

## 6. 全局参数
- `--token <token>`：显式传入 Token（优先级最高）
- `--api-base <url>`：自定义 API 地址（GHES）
- `--timeout <duration>`：HTTP 超时（默认 `15s`）
- `--json`：JSON 输出（支持的命令生效）
- `--verbose`：输出调试日志（不打印敏感信息）

## 7. 输出与错误码

### 7.1 输出约定
- 默认输出：面向人类可读
- `--json`：结构化输出，便于脚本集成

### 7.2 建议错误码
- `0`：成功
- `10`：认证失败（Token 缺失或无效）
- `11`：权限不足
- `12`：仓库或路径不存在
- `13`：参数错误
- `14`：网络或超时错误
- `15`：被限流
- `16`：本地文件写入失败
- `17`：用户取消操作

## 8. 常见使用流程
```bash
# 1) 配置 token
export GITHUB_TOKEN=ghp_xxx

# 2) 检查认证
ghrepo auth check

# 3) 查看目录
ghrepo ls owner/repo path/to/dir --ref main

# 4) 查看文件
ghrepo cat owner/repo path/to/file --ref main

# 5) 下载目录
ghrepo get owner/repo path/to/dir --out ./backup --ref main

# 6) 上传文件
ghrepo put owner/repo path/to/new-file.txt -m "add file" --file ./local.txt

# 7) 删除文件
ghrepo rm owner/repo path/to/old-file.txt -m "remove file"
```

## 9. 兼容性说明
- 支持 macOS / Linux
- 后续通过 Go 交叉编译支持 `amd64` 和 `arm64`
- Windows 支持可在 v0.2 评估

