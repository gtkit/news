// @Author xiaozhaofu 2022/10/31 20:13:00
package test

import (
	"testing"

	"gitlab.superjq.com/go-tools/news"
)

func TestFsWarnText(t *testing.T) {
	fsurl := ""
	news.FsWarnText(fsurl, "测试飞书消息包的文本消息")
}
