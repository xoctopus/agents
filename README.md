# agents

聚合 `xoctopus` 生态中的 Agent Skills, 作为统一的技能入口.

本仓库通过 Go module 依赖声明所需 skills, 再用 `genx/pkg/agent` 安装器把各模块提供的
skill 链接到本地 `.agents/skills`.

## 做了什么

1. **统一安装入口**: 在 `go.mod` 的直接依赖上标注 `// +skill:<name>`, 一次安装多个模块的
   skills.
2. **版本随依赖**: skill 内容来自对应 module 版本 (含 `replace`), 与代码依赖对齐.
3. **本地可发现**: 安装结果在 `.agents/skills/<name>`, 供 Cursor / Agent 读取
   `SKILL.md` 与 references.
4. **本仓库自带 skill**: `base` (文档与提交信息约定), 其余 skill 来自依赖模块.

当前已声明的 skills:

| Skill | 提供模块                    | 说明概要                        |
|-------|-----------------------------|---------------------------------|
| base  | 本仓库                      | 文档与提交信息规范              |
| appx  | `github.com/xoctopus/confx` | 组装 confx 可执行应用           |
| kg    | `github.com/xoctopus/confx` | 结构化 cache key 命名           |
| concx | `github.com/xoctopus/concx` | 受约束并发 (生命周期/编排/通信) |
| genx  | `github.com/xoctopus/genx`  | 代码生成扩展与 skills 安装      |
| logx  | `github.com/xoctopus/logx`  | 结构化日志与 span 上下文        |
| sqlx  | `github.com/xoctopus/sqlx`  | SQL 构建, 模型与生成接入        |
| testx | `github.com/xoctopus/x`     | 断言与 BDD 测试约定             |

安装后的布局示意:

```text
.agents/
  .gitignore          # 忽略已安装的 skills/<name>
  skills/
    base/             # 本仓库本地 skill
    appx -> <module>/.agents/skills/appx
    ...
```

约定细节见 `github.com/xoctopus/genx/pkg/agent` 与 genx skill
中的 [skills-installation.spec.md](https://github.com/xoctopus/genx/blob/main/.agents/skills/genx/references/skills-installation.spec.md).

## skill-install 命令

CLI 入口: `cmd/skill-install`. 在仓库根目录执行:

```bash
go run ./cmd/skill-install <command> [flags]
```

### install

安装 skills 并桥接到指定 Agent 目录. 流程:

1. (可选, `--upgrade`/`-u`) `GOWORK=off go get -u ./...` + `go mod tidy` 更新直接依赖
2. 按 `go.mod` 中 `// +skill:<name>` 安装到 `.agents/skills/<name>`
3. 将 `.agents/skills` 下每个 skill 强制链接到 Agent 目标目录

```bash
# 项目级, 安装全部 agent (claude / cursor / codex)
go run ./cmd/skill-install install

# 项目级, 仅链接到 .cursor/skills
go run ./cmd/skill-install install --name cursor

# 项目级, 仅链接到 .claude/skills
go run ./cmd/skill-install install --name claude

# 用户级, 链接到 ~/.cursor/skills
go run ./cmd/skill-install install --name cursor --mode 1

# 安装前先升级依赖到最新
go run ./cmd/skill-install install -u
```

| 参数               | 说明                                                                    |
|--------------------|-------------------------------------------------------------------------|
| `--name`           | 可选. `cursor` / `claude` / `codex`; 不指定则安装全部                   |
| `--path`           | 可选. 自定义目标目录, 需与 `--name` 同时使用; 不指定则用 Agent 默认路径 |
| `--mode`           | `0` 项目级(默认), `1` 当前用户 (`~` 下)                                 |
| `-u` / `--upgrade` | 可选. 安装前执行 `go get -u` 更新依赖                                   |

Agent 默认目标路径:

| Agent  | 项目级 (`mode=0`) | 用户级 (`mode=1`)  |
|--------|-------------------|--------------------|
| cursor | `.cursor/skills`  | `~/.cursor/skills` |
| claude | `.claude/skills`  | `~/.claude/skills` |
| codex  | `.agents/skills`  | `~/.agents/skills` |

`codex` 项目级时源与目标相同 (均为 `.agents/skills`), 跳过二次链接, 避免自引用 symlink.

安装后 Cursor 布局示例:

```text
.agents/skills/concx/         # skill 源 (指向 module cache)
.cursor/skills/concx -> .../.agents/skills/concx
```

### version

打印 `go.mod` 中带 `// +skill:` 标注的直接依赖及对应 skill 名称:

```bash
go run ./cmd/skill-install version
```

输出示例:

```text
github.com/xoctopus/concx@v0.2.3
skills:
	concx
github.com/xoctopus/confx@v0.6.0
skills:
	appx
	kg
```

## 如何安装 skills

### 前置条件

- 在仓库根目录执行 (需能读取 `go.mod`)
- 各 skill 提供方在模块内提供 `.agents/skills/<name>/` (至少含 `SKILL.md`)

### 安装

```bash
# 安装全部 agent
go run ./cmd/skill-install install

# 仅安装 cursor
go run ./cmd/skill-install install --name cursor
```

`genx/pkg/agent` 安装阶段会:

1. 解析当前 `go.mod` 中带 `// +skill:<name>` 的直接依赖
2. 在 `GOMODCACHE` (或 `replace` 路径) 定位 `<module>/.agents/skills/<name>`
3. 在 `.agents/skills/<name>` 创建 symlink
4. 向 `.agents/.gitignore` 追加 `skills/<name>` (幂等)

### 新增 / 变更 skill

1. 在 `go.mod` 对应 `require` 上方增加 `// +skill:<name>` (一个依赖可标注多个)
2. 确保该模块已发布 (或 `replace` 到本地) 且存在 `.agents/skills/<name>/`
3. 如需 blank import 固定依赖, 可在 `cmd/skill-install/main.go` 中补充
4. 重新执行 `install` 命令

示例:

```go
require (
	// +skill:logx
	github.com/xoctopus/logx v0.3.9
)
```

本地开发可用 `replace` 指向兄弟目录, 安装器会优先使用 replace 路径中的 skill 源.

### 排查

| 现象                      | 检查项                                                      |
|---------------------------|-------------------------------------------------------------|
| 安装报错找不到 skill 目录 | 模块版本是否含 `.agents/skills/<name>`, 或 replace 是否生效 |
| symlink 指向旧内容        | 重新执行 `install`, 会强制覆盖目标目录中的链接              |
| 未安装某个 skill          | `go.mod` 是否写了 `// +skill:<name>`, 且为直接依赖          |
| Agent 读不到 skill        | 是否已执行 `install`（或 `install --name <agent>`）         |
