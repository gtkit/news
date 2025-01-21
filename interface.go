package news

import (
	"context"
)

type (
	// AppNews 应用消息推送接口
	AppNewsInterface interface {
		Expired() int
		SendImageMsg(ctx context.Context, openID, filepath string) error // 发送图片消息
		SendTextMsg(ctx context.Context, openID, text string) error      // 发送文本消息
	}

	// AppCacher 应用缓存接口
	AppCacher interface {
		Get() AppNewsInterface
		Set(value AppNewsInterface)
	}
)
