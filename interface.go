package news

import (
	"context"
)

type (
	// AppNews 应用消息推送接口
	AppNewser interface {
		Expired() int
		SendImageMsg(ctx context.Context, openID, filepath string) error // 发送图片消息
		SendTextMsg(ctx context.Context, openID, text string) error      // 发送文本消息
	}

	// AppCacher 应用缓存接口
	AppCacher interface {
		Get() AppNewser
		Set(value AppNewser)
	}
)

// 发送图片消息
func SendImageMsg(ctx context.Context, a AppNewser, openID, filepath string) error {
	return a.SendImageMsg(ctx, openID, filepath)
}

// 发送文本消息
func SendTextMsg(ctx context.Context, a AppNewser, openID, msg string) error {
	return a.SendTextMsg(ctx, openID, msg)
}
