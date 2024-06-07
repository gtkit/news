package test

import (
	"testing"

	"github.com/gtkit/news"
)

func TestFsWarnText(t *testing.T) {
	fsurl := "https://open.feishu.cn/open-apis/bot/v2/hook/08498bc0-11e5-4b21-9a63-83fa5ceab488"
	news.FsNew(fsurl).RichText("我的测试标题", "测试飞书消息包的富文本消息", "111111111")
}
