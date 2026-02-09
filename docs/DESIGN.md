# ghrepo 设计文档（基于 USAGE v0.1）

## 1. 文档关系与范围
本设计文档基于 `/Users/jsonlee/Projects/githubRAGCli/docs/USAGE.md` 的命令契约，覆盖 v0.1 的只读能力：
- `auth check`
- `ls`
- `cat`
- `get`
- `stat`

不包含写操作（提交、删除、创建 PR 等）。

## 2. 设计目标
- 与使用文档保持一致的 CLI 体验
- 对私有仓库可用（基于 Token）
- 默认输出清晰，`--json` 输出稳定
- 错误码可被脚本可靠判断
- 为后续 Homebrew 发布保留标准化构建路径

## 3. 技术选型
- 语言：Go（1.22+）
- CLI 框架：`cobra`
- HTTP：标准库 `net/http`
- JSON 处理：标准库 `encoding/json`
- 构建发布：`goreleaser`（v0.2 接入）

说明：v0.1 采用直接 REST 封装，避免额外 SDK 抽象成本，便于控制错误映射与输出格式。

## 4. 模块划分

```text
cmd/ghrepo/main.go                # 入口
internal/cli/root.go              # 全局参数与命令注册
internal/cli/cmd_auth.go          # auth check
internal/cli/cmd_ls.go            # ls
internal/cli/cmd_cat.go           # cat
internal/cli/cmd_get.go           # get
internal/cli/cmd_stat.go          # stat
internal/config/config.go         # token/api-base/timeout 解析
internal/githubapi/client.go      # GitHub API 请求封装
internal/service/repo_service.go  # 业务逻辑（列目录、下载、类型判断）
internal/output/print.go          # text/json 输出统一
internal/errors/codes.go          # 退出码与错误类型
```

## 5. 命令到 API 映射

### 5.1 `auth check`
- API：`GET /user`
- 目的：验证 Token 有效性，并返回用户信息

### 5.2 `ls`
- 非递归：`GET /repos/{owner}/{repo}/contents/{path}?ref={ref}`
- 递归：先用 `contents` 获取目录 SHA，再调用
  `GET /repos/{owner}/{repo}/git/trees/{sha}?recursive=1`

### 5.3 `cat`
- API：`GET /repos/{owner}/{repo}/contents/{path}?ref={ref}`
- 文件类型判定后输出内容（base64 解码或 raw 媒体类型）

### 5.4 `get`
- 文件下载：复用 `cat` 文件读取逻辑后写本地
- 目录下载：先取目录树，再遍历文件逐个下载并写入本地

### 5.5 `stat`
- API：`GET /repos/{owner}/{repo}/contents/{path}?ref={ref}`
- 输出 `type/path/sha/size/download_url`

## 6. 配置与参数解析

### 6.1 Token 解析顺序
1. `--token`
2. `GITHUB_TOKEN`
3. `GH_TOKEN`

### 6.2 其他全局参数
- `--api-base`：默认 `https://api.github.com`
- `--timeout`：默认 `15s`
- `--json`
- `--verbose`

### 6.3 仓库参数解析
- 输入格式：`owner/repo`
- 解析失败即参数错误（错误码 `13`）

## 7. 数据模型（核心）

```go
type RepoRef struct {
  Owner string
  Repo  string
  Ref   string
}

type Entry struct {
  Type        string // file | dir
  Path        string
  Sha         string
  Size        int64
  DownloadURL string
}
```

## 8. 错误处理与退出码
与使用文档对齐：
- `0` 成功
- `10` 认证失败
- `11` 权限不足
- `12` 仓库或路径不存在
- `13` 参数错误
- `14` 网络/超时
- `15` 限流
- `16` 本地写入失败

映射策略：
- `401` -> `10`
- `403` 且权限不足 -> `11`
- `404` -> `12`
- `403` 且限流头触发 -> `15`
- 本地 I/O 异常 -> `16`

## 9. 下载与文件写入策略（`get`）
- 文件：
  - 若 `--out` 目标为目录，则使用原文件名
  - 若 `--out` 为文件路径，按该路径写入
- 目录：
  - 按仓库相对路径写入 `--out`
  - 自动创建父目录
- 覆盖策略：
  - 默认不覆盖；存在即报错
  - `--overwrite` 启用覆盖

## 10. 输出设计
- 文本输出：面向终端阅读，字段精简
- JSON 输出：稳定字段名，便于脚本解析
- `--verbose`：输出请求路径与重试信息，但不输出 token

## 11. 非功能设计
- 超时：全局 HTTP client timeout
- 重试：v0.1 仅对网络抖动进行有限重试（指数退避）
- 速率限制：识别 `X-RateLimit-Remaining` 与 `X-RateLimit-Reset`
- 安全：日志脱敏，避免 token 泄露

## 12. 测试策略
- 单元测试：
  - 参数解析
  - 错误码映射
  - 输出格式（text/json）
- 集成测试：
  - 使用 `httptest` 模拟 GitHub API
  - 覆盖 `ls/cat/get/stat` 关键路径
- 手工测试：
  - 公开仓库 + 私有仓库 + GHES endpoint（如可用）

## 13. 发布与 Homebrew 方案（v0.2）
- 使用 GitHub Actions 构建多平台产物
- 使用 `goreleaser` 生成 release 与 formula
- 发布 tap 仓库：`<yourname>/homebrew-tap`
- 用户安装：
  - `brew tap <yourname>/tap`
  - `brew install ghrepo`

## 14. 迭代计划
1. v0.1：核心命令与错误码落地
2. v0.2：Brew 发布、并发下载、`include/exclude` 过滤
3. v0.3：缓存、更完善限流与重试策略

