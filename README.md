# monitor-lite-api

IPTV 直播源监控与管理服务（Web 管理端名称：wtv-manage）。提供直播源的采集、检测、分组、订阅、EPG 及 M3U 播放列表输出，内置 Web 管理界面与若干电视台代理接口。

## 功能概览

### 媒体管理

- 维护本地 SQLite 中的直播源（增删改查、批量操作）
- 通过 FFmpeg 检测源可用性，记录分辨率、码率、失败次数等
- 按规则引擎自动分组、频道名标准化 from `etc/rule.yaml`
- 支持黑名单过滤、失效源自动删除

### 订阅池

- 添加 M3U/M3U8 订阅地址，自动解析并入库
- 支持手动抓取与定时自动读取（可在系统设置中配置每日时间）

### 播放列表输出

- `/v1/tv/super` — 完整 M3U 播放列表（国内过滤）
- `/v1/tv/w/:plusKey` — DIYP 格式输出，`plusKey` 在系统设置中配置
- `/cus/:path` — 自建输出策略，可按码率、分组等条件定制
- `/logo/*path` — 频道 Logo 静态资源

### 代理接口

| 路径 | 说明 |
|------|------|
| `/v1/proxy/fengshows?t=zx` | 凤凰资讯 |
| `/v1/proxy/fengshows?t=zw` | 凤凰中文 |
| `/v1/proxy/fengshows?t=xg` | 凤凰香港台 |
| `/v1/proxy/ptbtv?t=pt1` | 莆田一套 |
| `/v1/proxy/ptbtv?t=pt2` | 莆田二套 |
| `/v1/proxy/ptbtv?t=xy` | 仙游台 |
| `/v1/proxy/iqilu` | 齐鲁网直播源代理 |

### EPG 节目单

- `/v1/epg/epgList` — 查询节目单
- `/v1/epg/diyp` — DIYP 格式 EPG
- `/v1/epg/collect` — 手动触发采集

### Web 管理界面

访问 `http://<host>:<port>/admin`，主要模块：

| 模块 | 说明 |
|------|------|
| 媒体维护 | 频道列表管理、检测、分组 |
| 监控画面 | WebSocket + FFmpeg 在线预览 |
| 监控引擎 | 编辑 `rule.yaml` 分组与频道映射规则 |
| 订阅池 | 管理 M3U 订阅源 |
| 系统设置 | 自动守护、自动分组、EPG、黑名单、DIYP 链接等 |

## 技术栈

- **后端**：Go 1.23 + Gin + GORM + SQLite
- **前端**：Vue 3 + Element Plus + Vite（嵌入二进制 `/admin` 路径）
- **依赖**：FFmpeg（启动前必须已安装）

## 配置

首次启动会自动从内置模板生成配置文件：

| 文件 | 说明 |
|------|------|
| `etc/tv.yaml` | 服务端口、超时、管理员账号等 |
| `etc/rule.yaml` | 频道分组与名称映射规则 |
| `etc/sqlite.db` | 数据存储 |

`etc/tv.yaml` 示例：

```yaml
Name: tv
Host: 0.0.0.0
Port: 9876
Timeout: 50000
Username: "admin"
Password: "admin123"
PlusKey: "diyp"
```

## 本地开发

### 前置条件

- Go 1.23+
- Node.js 16+
- FFmpeg
- 在 monorepo 根目录 `ygczpro-api` 下执行（依赖 `common` 模块）

### 构建前端

```bash
cd ui
npm install
npm run build-only
```

### 启动后端

```bash
# 在 monitor-lite-api 目录下
go run .
# 或指定配置文件
go run . -f etc/tv.yaml
```

默认管理地址：`http://localhost:9876/admin`

## 后台定时任务

| 周期 | 任务 | 开关 |
|------|------|------|
| 启动时 | 抓取订阅、更新分组 | — |
| 每 10 分钟 | 自动检测直播源 | 系统设置「自动守护」 |
| 每 15 分钟 | 从源池自动搜源 | 系统设置「自动搜媒体」 |
| 每 120 分钟 | 按规则重新分组 | 系统设置「自动分组」 |
| 每天 08:00 | 采集 EPG | 系统设置「自动更新 EPG」 |

## 主要 API

```
GET  /v1/tv/json          搜索/查询频道
GET  /v1/tv/page          分页列表
POST /v1/tv/update        更新频道
POST /v1/tv/check         检测单个源
POST /v1/tv/checkAll      批量检测
POST /v1/tv/batchupdate   批量更新
POST /v1/tv/batchdelete   批量删除
GET  /v1/tv/rule/get      获取规则
POST /v1/tv/rule/update   更新规则
GET  /v1/setting/find     获取系统设置
POST /v1/setting/update   更新系统设置
GET  /v1/video/play       WebSocket 视频预览（参数 url）
GET  /v1/subscriber/*     订阅池管理
GET  /v1/selfout/*        自建输出策略
```

除公开路径（播放列表、代理、Logo、登录等）外，API 需在请求头携带 `session-id`（JWT Token）。

## 部署

### 自动发布

推送版本标签即可触发 GitHub Actions 自动打包并创建 Release（支持 Linux / Windows / macOS，amd64 与 arm64）：

```bash
git tag v1.5.0
git push origin v1.5.0
```

也可在 GitHub Actions 页面手动触发 `monitor-lite-api Release` 工作流。

Release 产物命名规则：

| 平台 | 文件名 |
|------|--------|
| Linux | `monitor-lite-api-linux-amd64` / `monitor-lite-api-linux-arm64` |
| Windows | `monitor-lite-api-windows-amd64.exe` / `monitor-lite-api-windows-arm64.exe` |
| macOS | `monitor-lite-api-darwin-amd64` / `monitor-lite-api-darwin-arm64` |

### 一键安装

将 `<owner>/<repo>` 替换为实际的 GitHub 仓库（默认 `bsxbl/ygczpro-api`）。

**Linux（需 root）**

```bash
# 安装
curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install.sh" | sudo bash -s install

# 更新
curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install.sh" | sudo bash -s update

# 卸载
curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install.sh" | sudo bash -s uninstall
```

**macOS**

```bash
curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install-macos.sh" | bash -s install
```

**Windows（管理员 PowerShell）**

```powershell
irm "https://github.com/<owner>/<repo>/releases/latest/download/install.ps1" | iex
# 更新: irm ... | iex -Update
# 卸载: irm ... | iex -Uninstall
```

可通过环境变量自定义仓库：`GITHUB_REPO=owner/repo`。

默认安装路径：Linux `/opt/wtv`，macOS `~/.wtv`，Windows `C:\Program Files\wtv`。服务监听 `9876` 端口。

### Docker 镜像

```bash
docker tag ea7b7ca1cdd9 ygcz/wtv-server:1.2.5
```

参考：

- https://blog.51cto.com/u_14850/11271288
- https://blog.csdn.net/m624197265/article/details/141719515

### 交叉编译

在 monorepo 根目录执行 `./build.sh`，选择 `monitor-lite-api` 即可同时构建前端并打包 `linux/amd64`、`linux/arm64` 二进制。
