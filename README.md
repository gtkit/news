### news

#### 飞书webhook消息发送
```go
fsurl := "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"

fs.NewWebHook(fsurl).Send(消息内容) // 纯文本消息

fs.NewWebHook(fsurl).Send("我的测试标题", "测试飞书消息包的富文本消息") // 富文本消息, 标题和内容参数可自定义
```
#### 应用消息发送
```go
package main

import (
	"context"
	"github.com/gtkit/news/fs"
)

func main() {
	app_id := "xxx"
	app_secret := "xxx"
        ctx := context.Background()
	app, _ := fs.NewInternalApp(ctx, app_id, app_secret)
	
	shotname := "test.png"
	// 上传图片
	img, _ := app.UploadImage(ctx, shotname)
	// 发送图片消息
	_ := app.SendImageMsg(ctx, openid, img.ImageKey())
	// 发送文本消息
	_ := app.SendTextMsg(ctx, openid, "测试应用消息包的文本消息")
	// app token 过期时间获取
	expire_time := app.Expired()
}

```
