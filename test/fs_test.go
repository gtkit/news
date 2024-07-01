package test

import (
	"testing"

	"github.com/gtkit/news"
)

func TestFsWarnText(t *testing.T) {
	fsurl := "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"
	news.FsNew(fsurl).Send("我的标题1")
	news.FsNew(fsurl).Send("我的标题2", "我的内容")
}
