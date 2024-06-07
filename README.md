### news

#### 飞书消息发送
```
fsurl := "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"

news.FsNew(fsurl).Send(消息内容) // 纯文本消息

news.FsNew(fsurl).Send("我的测试标题", "测试飞书消息包的富文本消息") // 富文本消息, 标题和内容参数可自定义
```
