package fs

import (
	"github.com/gtkit/news"
)

// 检查接口实现
func _() {
	var (
		_ news.AppNewser = (*InternalApp)(nil)
	)
}
