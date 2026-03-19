# news

[![Go Reference](https://pkg.go.dev/badge/github.com/gtkit/news.svg)](https://pkg.go.dev/github.com/gtkit/news)
[![Go Version](https://img.shields.io/badge/go-1.26+-blue.svg)](https://go.dev)

多平台 IM 机器人消息推送 Go 库，支持**飞书 (Lark)**、**企业微信 (WeCom)**、**钉钉 (DingTalk)**。

本仓库包含两个独立模块：

| 模块 | 导入路径 | 最新版本 | 说明 |
|-----|---------|---------|------|
| **v1** | `github.com/gtkit/news` | v1.1.9 | 飞书专用：应用消息 + Webhook |
| **v2** | `github.com/gtkit/news/v2` | v2.0.0 | 三平台统一 Provider 接口 + Gin 集成 |

---

## 目录

- [v1 — 飞书专用](#v1--飞书专用)
  - [安装 v1](#安装-v1)
  - [Webhook 消息](#webhook-消息)
  - [应用消息（需 access token）](#应用消息需-access-token)
- [v2 — 多平台 Provider](#v2--多平台-provider)
  - [安装 v2](#安装-v2)
  - [快速开始](#快速开始)
  - [飞书](#飞书-feishu)
  - [企业微信](#企业微信-wecom)
  - [钉钉](#钉钉-dingtalk)
  - [图片上传](#图片上传)
  - [多平台广播](#多平台广播)
  - [Gin 框架集成](#gin-框架集成)
  - [消息类型支持矩阵](#消息类型支持矩阵)
- [项目结构](#项目结构)
- [版本发布](#版本发布)
- [设计原则](#设计原则)
- [从 v1 迁移到 v2](#从-v1-迁移到-v2)

---

## v1 — 飞书专用

v1 模块专注于飞书平台，提供 **Webhook 群消息** 和 **应用消息**（需 `app_id` / `app_secret`）两种能力。

### 安装 v1

```bash
go get github.com/gtkit/news
```

### Webhook 消息

无需 `access_token`，拿到飞书群机器人的 Webhook 地址即可发送：

```go
import "github.com/gtkit/news/fs"

// 发送纯文本消息
err := fs.WebHookSend(fsurl, "服务上线通知：order-service v1.0.0")

// 发送富文本消息（标题 + 正文 + @所有人）
err := fs.WebHookSend(fsurl, "我的标题", "我的内容")

// 发送富文本消息（标题 + 正文 + 超链接提示 + 超链接 + @所有人）
err := fs.WebHookSend(fsurl, "标题", "正文内容", "点击查看", "https://example.com")

// 带 context 的版本（推荐用于 HTTP handler 中）
err := fs.WebHookSendCtx(ctx, fsurl, "hello world")
```

### 应用消息（需 access token）

通过飞书自建应用的 `app_id` / `app_secret` 获取 token，可以向指定用户发送消息：

```go
import (
    "context"
    "github.com/gtkit/news/fs"
)

ctx := context.Background()

// 创建应用实例（自动获取 tenant_access_token）
app, err := fs.NewInternalApp(ctx, "cli_xxx", "secret_xxx")

// 发送文本消息
err = app.SendTextMsg(ctx, "ou_xxx", "你好，这是一条应用消息")

// 上传图片（独立调用，获取 image_key）
uploadResp, err := app.UploadImage(ctx, "/path/to/image.png")
imageKey := uploadResp.ImageKey()

// 发送图片消息（内部自动上传 + 发送）
err = app.SendImageMsg(ctx, "ou_xxx", "/path/to/image.png")

// 下载图片
err = app.DownloadImage(ctx, imageKey, "/path/to/save.png")

// 获取 token 过期时间（秒）
expire := app.Expired()
```

#### 带缓存

```go
type MyCache struct{}

func (c *MyCache) Get() news.AppNewser { return fs.EmptyApp() }
func (c *MyCache) Set(app news.AppNewser) { /* 存入 Redis 等 */ }

app, err := fs.NewInternalApp(ctx, appID, appSecret, &MyCache{})
```

---

## v2 — 多平台 Provider

v2 是完全重写的版本，提供统一的 `Provider` 接口兼容三个平台，支持 Gin 集成、并发广播、图片上传。

### 安装 v2

```bash
go get github.com/gtkit/news/v2
```

> 要求 Go 1.26+

### 快速开始

```go
import (
    "context"
    "github.com/gtkit/news/v2"
    "github.com/gtkit/news/v2/feishu"
)

fs, err := feishu.New(news.Config{
    WebhookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/your-token",
    Secret:     "your-secret", // 可选：签名校验
})

ctx := context.Background()
err = fs.SendText(ctx, "Hello from news/v2!")
```

### 飞书 (Feishu)

```go
import (
    "github.com/gtkit/news/v2"
    "github.com/gtkit/news/v2/feishu"
)

fs, _ := feishu.New(news.Config{
    WebhookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
    Secret:     "feishu-secret",
})

// 文本
fs.SendText(ctx, "hello", news.WithAtAll())
fs.SendText(ctx, "请审批", news.WithAtUsers("ou_xxx1", "ou_xxx2"))

// Markdown（自动封装为交互卡片）
fs.SendMarkdown(ctx, "部署报告", "**服务**: order\n**状态**: ✅")

// 富文本（飞书原生 post 格式）
msg := feishu.BuildRichText(
    "监控告警",
    "服务 order-service 响应超时。",
    &news.RichTextTag{Tag: "a", Text: "查看详情", Href: "https://grafana.example.com"},
    true,
)
fs.SendRichText(ctx, msg)

// 多行富文本
msg := feishu.BuildRichTextLines("发布通知",
    []news.RichTextTag{
        {Tag: "text", Text: "user-service 已部署。"},
    },
    []news.RichTextTag{
        {Tag: "text", Text: "版本: "},
        {Tag: "a", Text: "v2.1.0", Href: "https://github.com/xxx"},
    },
    []news.RichTextTag{
        {Tag: "at", UserID: "all"},
    },
)
fs.SendRichText(ctx, msg)

// 图片（已有 image_key）
fs.SendImage(ctx, &news.ImageMessage{ImageKey: "img_v2_xxx"})
```

### 企业微信 (WeCom)

```go
import (
    "github.com/gtkit/news/v2"
    "github.com/gtkit/news/v2/wecom"
)

wc, _ := wecom.New(news.Config{
    WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx",
})

// 文本
wc.SendText(ctx, "hello", news.WithAtAll())
wc.SendText(ctx, "请处理", news.WithAtUsers("zhangsan", "lisi"))

// Markdown（支持 <font color="warning">橙色</font> 等）
wc.SendMarkdown(ctx, "告警",
    `新增反馈<font color="warning">132例</font>，请注意。`)

// 图片（base64 + md5，自动计算）
wc.SendImageFromFile(ctx, "/path/to/alert.png")
```

### 钉钉 (DingTalk)

```go
import (
    "github.com/gtkit/news/v2"
    "github.com/gtkit/news/v2/dingtalk"
)

dt, _ := dingtalk.New(news.Config{
    WebhookURL: "https://oapi.dingtalk.com/robot/send?access_token=xxx",
    Secret:     "SEC-xxx",
})

// 文本
dt.SendText(ctx, "hello", news.WithAtUsers("13800001111"))

// Markdown
dt.SendMarkdown(ctx, "杭州天气", "#### 杭州天气\n> 9度，西北风1级")

// 图片（公网 URL 嵌入 markdown）
dt.SendImageFromURL(ctx, "https://example.com/screenshot.png")

// Link 消息
dt.SendLink(ctx, "标题", "正文", "https://example.com", "https://example.com/pic.png")

// ActionCard（独立跳转按钮）
dt.SendActionCard(ctx, &dingtalk.ActionCard{
    Title: "CI 构建失败",
    Text:  "### 构建失败\n分支: main",
    Buttons: []dingtalk.Button{
        {Title: "查看日志", ActionURL: "https://ci.example.com/123"},
        {Title: "重新构建", ActionURL: "https://ci.example.com/123/retry"},
    },
})

// FeedCard
dt.SendFeedCard(ctx, []dingtalk.FeedLink{
    {Title: "PR #123 已合并", MessageURL: "https://github.com/...", PicURL: "https://..."},
})
```

### 图片上传

三个平台的图片处理方式不同，v2 针对各自 API 特性提供了最合适的封装：

#### 飞书 — 独立上传 API（需 tenant_access_token）

```go
import "github.com/gtkit/news/v2/feishu"

// 方式一：独立上传，获取 image_key 后自行使用
resp, err := feishu.UploadImageFromFile(ctx, "t-tenant_access_token", "/path/to/img.png")
imageKey := resp.ImageKey()

// 方式二：从 io.Reader 上传（适用于内存中的图片、HTTP 响应体等）
resp, err := feishu.UploadImageFromReader(ctx, "t-token", "screenshot.png", imageReader)

// 方式三：一步完成「上传 + 通过 webhook 发送」
err := fsProvider.SendImageFromFile(ctx, "t-token", "/path/to/img.png")
```

> `tenant_access_token` 由调用方通过 `fs.NewInternalApp` 或其他方式获取。

#### 企业微信 — 无需上传，自动 base64 + md5

```go
import "github.com/gtkit/news/v2/wecom"

// 方式一：直接从文件发送（内部自动读取 → base64 → md5）
err := wcProvider.SendImageFromFile(ctx, "/path/to/img.png")

// 方式二：独立构建 ImageMessage（适用于已有 []byte 的场景）
img := wecom.BuildImageMessage(imageBytes)
err := wcProvider.SendImage(ctx, img)
```

#### 钉钉 — 只支持公网图片 URL

```go
import "github.com/gtkit/news/v2/dingtalk"

// 通过公网可访问的图片 URL 发送（内部包装为 markdown）
err := dtProvider.SendImageFromURL(ctx, "https://example.com/screenshot.png")
```

### 多平台广播

```go
fs, _ := feishu.New(news.Config{WebhookURL: "..."})
wc, _ := wecom.New(news.Config{WebhookURL: "..."})
dt, _ := dingtalk.New(news.Config{WebhookURL: "..."})

multi, _ := news.NewMulti(fs, wc, dt)

// 并发发送到三个平台（内部使用 wg.Go()）
err := multi.SendText(ctx, "全平台通知", news.WithAtAll())
err := multi.SendMarkdown(ctx, "标题", "**正文**")
```

### Gin 框架集成

```go
import (
    "github.com/gtkit/news/v2"
    "github.com/gtkit/news/v2/ginews"
)

mgr := news.NewManager(fs, wc, dt)

r := gin.Default()
r.Use(ginews.Middleware(mgr))

r.POST("/alert", func(c *gin.Context) {
    mgr := ginews.MustFrom(c)

    // 发到指定平台
    mgr.Feishu().SendText(c.Request.Context(), "Gin handler 告警")

    // 广播所有平台
    multi, _ := mgr.Multi()
    multi.SendText(c.Request.Context(), "全平台告警", news.WithAtAll())

    c.JSON(200, gin.H{"status": "ok"})
})
```

**ginews API 速查：**

| 函数 | 说明 |
|-----|------|
| `ginews.Middleware(mgr)` | 注入 Manager 到 Gin context |
| `ginews.From(c)` | 获取 Manager，未找到返回 nil |
| `ginews.MustFrom(c)` | 获取 Manager，未找到 panic |
| `ginews.ProviderFrom(c, platform)` | 获取指定平台 Provider |

### 消息类型支持矩阵

| 消息类型 | 飞书 | 企业微信 | 钉钉 |
|---------|------|---------|------|
| 文本 | ✅ | ✅ | ✅ |
| Markdown | ✅ 交互卡片 | ✅ | ✅ |
| 富文本 (Post) | ✅ 原生 | ✅ → MD | ✅ → MD |
| 图片上传 | ✅ `UploadImageFromFile` | ✅ `SendImageFromFile` (auto base64) | ❌ 无上传 API |
| 图片发送 | ✅ image_key | ✅ base64+md5 | ✅ picURL→MD |
| Link | — | — | ✅ |
| ActionCard | — | — | ✅ |
| FeedCard | — | — | ✅ |

### @用户

| 平台 | `WithAtAll()` | `WithAtUsers()` 参数 |
|------|-------------|---------------------|
| 飞书 | ✅ | user_id 或 open_id |
| 企业微信 | ✅ | userid |
| 钉钉 | ✅ | 手机号 |

---

## 项目结构

```
github.com/gtkit/news
├── version.go                  v1 版本号: v1.1.9
├── interface.go                AppNewser / AppCacher 接口
├── go.mod                      module github.com/gtkit/news
├── Makefile                    make tag (v1) / make v2tag (v2)
│
├── fs/                         v1 飞书实现
│   ├── app.go                  NewInternalApp / EmptyApp / Expired
│   ├── model.go                飞书官方响应结构体
│   ├── message.go              SendImageMsg / SendTextMsg
│   ├── image.go                UploadImage / DownloadImage
│   └── webhook.go              WebHookSend / WebHookSendCtx
│
└── v2/                         module github.com/gtkit/news/v2
    ├── version.go              v2 版本号: v2.0.0
    ├── go.mod                  独立 go.mod
    ├── provider.go             Provider 接口 / Config / Stats
    ├── richtext.go             RichTextToMarkdown
    ├── manager.go              Manager + Multi 并发广播
    ├── feishu/
    │   ├── webhook.go          text / markdown / richtext / image
    │   └── image.go            UploadImageFromFile / UploadImageFromReader / SendImageFromFile
    ├── wecom/
    │   ├── webhook.go          text / markdown / image
    │   └── image.go            SendImageFromFile / BuildImageMessage
    ├── dingtalk/
    │   ├── webhook.go          text / markdown / link / ActionCard / FeedCard
    │   └── image.go            SendImageFromURL
    ├── ginews/
    │   └── middleware.go       Gin 中间件
    └── internal/
        ├── http.go             PostJSON（带 LimitReader）
        └── sign.go             FeishuSign / DingTalkSignedURL
```

## 版本发布

v1 和 v2 使用独立的 `version.go` 管理版本号，通过 `Makefile` 自动递增 patch 版本并打 tag：

```bash
# v1 发版：version.go v1.1.9 → v1.1.10，打 tag v1.1.10
make tag

# v2 发版：v2/version.go v2.0.0 → v2.0.1，打 tag v2/v2.0.1
make v2tag

# 查看最近 tag
make gittag
```

## 设计原则

| 原则 | 实现 |
|-----|------|
| **不可变构造** | 所有 Webhook 在 `New()` 后字段不再修改，天然并发安全 |
| **无锁设计** | Stats 使用 `atomic.Int64`，Manager.defaults 使用 `atomic.Value` |
| **HTTP 连接复用** | `Config.Freeze()` 在构造时冻结 `http.Client`，热路径零分配 |
| **接口隔离** | `Provider` 接口只定义通用能力，平台专属方法通过具体类型暴露 |
| **优雅降级** | `SendRichText` 在企业微信/钉钉自动转为 Markdown |
| **标准 error** | 全部 `fmt.Errorf("%w")`，支持 `errors.Is` / `errors.As` |
| **响应体限制** | `io.LimitReader(1MB)` 防止 OOM |
| **Go 1.25+ 特性** | `wg.Go()` 并发广播，`omitzero` JSON tag |
| **零外部依赖** | v2 核心仅依赖标准库；ginews 依赖 Gin；v1 依赖 gtkit/json |

## 从 v1 迁移到 v2

| v1 (github.com/gtkit/news) | v2 (github.com/gtkit/news/v2) |
|---|---|
| 只支持飞书 | 飞书 + 企业微信 + 钉钉 |
| `fs.WebHookSend(url, args...)` | `provider.SendText(ctx, text, opts...)` |
| 返回 `error`（优化后） | 返回 `error` |
| `fs.NewInternalApp` + `SendImageMsg` | `feishu.UploadImageFromFile` + `SendImage` |
| 无 Gin 集成 | `ginews.Middleware` + `ginews.MustFrom` |
| 无多平台支持 | `news.NewMulti` 并发广播 |
| 无签名支持 | 飞书 + 钉钉 HMAC-SHA256 签名 |

> v1 和 v2 可以在同一项目中共存，互不影响。v1 的飞书应用消息能力（`NewInternalApp`）在 v2 中没有对应物，如需该功能请继续使用 v1 的 `fs` 包。

## License

MIT
