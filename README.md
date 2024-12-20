### news

#### 飞书webhook消息发送
```go
fsurl := "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"

fs.WebHookNews(fsurl, "测试飞书消息包的纯文本消息内容") // 纯文本消息

fs.WebHookNews(fsurl, "我的测试标题", "测试飞书消息包的富文本消息") // 富文本消息, 标题和内容参数可自定义
```
#### 应用消息发送
```go
package main

import (
	"context"
	"github.com/gtkit/news/fs"
)

type Cache struct {
	key string
}

func main() {
	app_id := "xxx"
	app_secret := "xxx"
        ctx := context.Background()
        cacher := Cache{
            key: "xxx",
        }
	app, _ := fs.NewInternalApp(ctx, app_id, app_secret, cacher)
	
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

// Cache implements news.AppCacher
func (c *Cache) Get() news.AppNewsInterface {
        // implement your cache logic here
	
	app := fs.EmptyApp()
	
	return app
}
func (c *Cache) Set(app news.AppNewsInterface) {
	// implement your cache logic here
	return
}

```
