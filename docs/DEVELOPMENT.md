# ghrepo 开发文档（基于 DESIGN v0.1）

## 1. 目的与范围
本文档用于指导 `ghrepo` 的工程实现，目标是把 `/Users/jsonlee/Projects/githubRAGCli/docs/DESIGN.md` 转化为可执行开发任务。

当前范围（v0.1）：
- 只读命令：`auth check`、`ls`、`cat`、`get`、`stat`
- 支持 GitHub.com，预留 GHES（`--api-base`）
- 不包含任何写仓库能力

## 2. 开发环境

### 2.1 依赖版本
- Go: `1.22+`
- Git: 最新稳定版
- 可选：`golangci-lint`、`goreleaser`（v0.2 使用）

### 2.2 本地环境变量
```bash
export GITHUB_TOKEN=ghp_xxx
```

### 2.3 初始化命令（首次）
```bash
go mod init githubRAGCli
go get github.com/spf13/cobra@latest
```

## 3. 代码目录与职责
按设计文档落地以下结构：

```text
cmd/ghrepo/main.go
internal/cli/root.go
internal/cli/cmd_auth.go
internal/cli/cmd_ls.go
internal/cli/cmd_cat.go
internal/cli/cmd_get.go
internal/cli/cmd_stat.go
internal/config/config.go
internal/githubapi/client.go
internal/service/repo_service.go
internal/output/print.go
internal/errors/codes.go
```

职责约束：
- `cli`：只做参数解析与调用，不直接拼 HTTP 请求
- `githubapi`：只做 API 访问，不做业务流程编排
- `service`：聚合 API，处理目录/文件判定、递归、下载策略
- `output`：统一 text/json 输出，不在命令中散落 `fmt.Printf`
- `errors`：统一退出码与错误类型映射

## 4. 开发顺序（建议）

1. 基础骨架
- 建立 `main/root`，接入全局参数：`--token --api-base --timeout --json --verbose`
- 增加统一配置读取与校验

2. 客户端层
- 完成 GitHub API 客户端、请求头注入、超时控制
- 实现响应错误分类：`401/403/404/429/5xx`

3. 服务层
- 实现 `GetContent`、`GetTree`、`ReadFile`、`DownloadPath` 等核心方法
- 建立 `Entry` 统一模型

4. 命令层
- 按 `auth -> stat -> ls -> cat -> get` 顺序接入
- 每个命令完成 text/json 两种输出

5. 质量与收尾
- 单元测试 + 集成测试
- README/USAGE 对齐校验
- 版本号与构建脚本

## 5. 命令实现说明

## 5.1 `auth check`
流程：
1. 读取 token（优先级：flag > `GITHUB_TOKEN` > `GH_TOKEN`）
2. 调 `GET /user`
3. 打印用户与限流信息

失败处理：
- token 缺失/无效 -> 退出码 `10`

## 5.2 `stat`
流程：
1. 解析 `owner/repo` + `path` + `ref`
2. 调 `contents` 接口
3. 输出 `type/path/sha/size/download_url`

失败处理：
- 路径不存在 -> `12`

## 5.3 `ls`
流程：
1. 非递归：直接 `contents`
2. 递归：目录 `sha` -> `git/trees?recursive=1`
3. 转成 `[]Entry` 输出

失败处理：
- 路径是文件但调用 `ls`：返回参数/语义错误，退出码 `13`

## 5.4 `cat`
流程：
1. 调 `contents`
2. 校验目标类型必须是 `file`
3. 解码内容并输出 stdout

失败处理：
- 目标是目录 -> `13`

## 5.5 `get`
流程：
1. 判定 `path` 类型（file/dir）
2. 文件：下载并写入 `--out`
3. 目录：递归列文件并逐个下载写入
4. 应用覆盖策略（默认不覆盖，`--overwrite` 覆盖）

失败处理：
- 本地写失败 -> `16`

## 6. API 与数据契约

### 6.1 内部核心结构
```go
type RepoRef struct {
    Owner string
    Repo  string
    Ref   string
}

type Entry struct {
    Type        string `json:"type"`
    Path        string `json:"path"`
    Sha         string `json:"sha"`
    Size        int64  `json:"size,omitempty"`
    DownloadURL string `json:"download_url,omitempty"`
}
```

### 6.2 输出稳定性要求
- JSON 字段名不随版本随意变动
- 文本输出可以增强，但不能改变命令语义

## 7. 错误码与异常处理

固定退出码：
- `0` 成功
- `10` 认证失败
- `11` 权限不足
- `12` 资源不存在
- `13` 参数错误
- `14` 网络/超时
- `15` 限流
- `16` 本地写入失败

实现要求：
- 所有命令错误最终转换为统一错误类型
- 主程序仅在一个出口设置 `os.Exit(code)`
- 日志和错误信息均不得泄露 token

## 8. 测试开发策略

### 8.1 单元测试
覆盖点：
- `owner/repo` 解析
- token 优先级读取
- 错误码映射
- 输出序列化（`--json`）

建议命令：
```bash
go test ./...
```

### 8.2 集成测试（`httptest`）
覆盖点：
- `auth check` 成功/401
- `ls` 普通目录/递归目录
- `cat` 文件与目录误用
- `get` 文件下载与目录下载
- `stat` 不存在路径

### 8.3 手工验收
最小验收清单：
1. 公开仓库执行所有命令
2. 私有仓库执行 `ls/cat/get`
3. `--json` 输出可被 `jq` 解析
4. 错误码可在 shell 中正确判断（`echo $?`）

## 9. 性能与可靠性
- HTTP 超时默认 `15s`
- 失败重试：仅网络抖动类错误，指数退避（短重试）
- 下载目录时按顺序实现（v0.1），并发下载放到 v0.2
- 限流时输出 reset 时间，返回 `15`

## 10. 安全要求
- 禁止打印 token（包括 verbose）
- 错误日志中的请求头必须脱敏
- 不把 token 写入配置文件或缓存文件

## 11. CI/CD 与发布（阶段性）

v0.1：
- CI：`go test ./...`
- 构建：`go build -o ghrepo ./cmd/ghrepo`

v0.2：
- 接入 `goreleaser`
- 生成多平台产物（darwin/linux; amd64/arm64）
- 发布 Homebrew tap formula

## 12. 任务拆解（可直接建 issue）

1. `feat(cli-root)`: 初始化 root 命令与全局参数
2. `feat(config)`: token/api-base/timeout 解析
3. `feat(api-client)`: GitHub REST 客户端与错误分类
4. `feat(auth-check)`: 认证检查命令
5. `feat(stat)`: 路径元信息命令
6. `feat(ls)`: 目录列出（含 recursive）
7. `feat(cat)`: 文件读取命令
8. `feat(get)`: 下载命令与覆盖策略
9. `test(core)`: 单测与集成测试
10. `chore(release)`: 构建脚本与版本流程

## 13. 完成定义（DoD）
满足以下条件视为 v0.1 完成：
1. 五个命令均可运行并符合 `USAGE.md`
2. 错误码与 `DESIGN.md` 完全一致
3. `go test ./...` 全部通过
4. 在公开仓库和私有仓库手工验收通过
5. 文档（USAGE/DESIGN/DEVELOPMENT）保持一致

