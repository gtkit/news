package news

import (
	"context"
)

// AppNews 应用消息推送接口
type AppNewsInterface interface {
	Expired() int
	UploadImage(ctx context.Context, path string) (ImageKeyer, error)
	DownloadImage(ctx context.Context, imageKey, path string) error
	SendImageMsg(ctx context.Context, openID, imageKey string) error // 发送图片消息
	SendTextMsg(ctx context.Context, openID, text string) error      // 发送文本消息
}

// UploadImageInfo 上传图片返回信息
type ImageKeyer interface {
	ImageKey() string
}

type AppCacher interface {
	Get() AppNewsInterface
	Set(value AppNewsInterface)
}
