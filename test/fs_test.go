package test

import (
	"context"
	"log"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"

	"github.com/gtkit/news/fs"
)

var (
	appID     = "xxx"
	appSecret = "xxx"
	srcimg    = "/xxx/800000797.jpg"
	dstpath   = "/xxx/"
)

func TestFsWarnText(t *testing.T) {
	fsurl := "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"
	fs.WebHookNews(fsurl, "我的标题1")
	fs.WebHookNews(fsurl, "我的标题2", "我的内容")
}

func TestAccessToken(t *testing.T) {
	// NewApp("cli_9f3dd38ac5bbd00e", "WaVHcgdg2n9slTh5y7AutbNqBogZhdWJ")
	app, _ := fs.NewInternalApp(context.Background(), appID, appSecret)
	t.Logf("app: %+v", app)
	// tenantAccessToken := app.TenantAccessToken()
	// t.Log("tenantAccessToken: ", tenantAccessToken)
}

func TestUploadImage(t *testing.T) {

	ctx := context.Background()

	app, err := fs.NewInternalApp(ctx, appID, appSecret)
	if err != nil {
		t.Error("NewInternalApp error: ", err)
	}
	res, err := app.UploadImage(ctx, srcimg)
	if err != nil {
		t.Error("UploadImage error: ", err)
	}
	t.Logf("res: %+v", res)
	if err := app.DownloadImage(ctx, res.ImageKey(), dstpath+filepath.Base(srcimg)); err != nil {
		t.Error("UploadImage error: ", err)
	}
}

func TestSendImageMsg(t *testing.T) {
	ctx := context.Background()
	app, err := fs.NewInternalApp(ctx, appID, appSecret)
	if err != nil {
		t.Error("NewInternalApp error: ", err)
	}

	img, err := app.UploadImage(ctx, srcimg)
	if err != nil {
		t.Error("UploadImage error: ", err)
	}
	t.Logf("res: %+v", img)

	if err := app.SendImageMsg(ctx, "ou_1e3fc242928fa853dd2ed13b1db60bd3", img.ImageKey()); err != nil {
		t.Error("SendImageMessage error: ", err)
		return
	}
	t.Log("SendImageMessage success")
}

func TestErr(t *testing.T) {
	err := errors.New("test error")
	log.Printf("access token post error:%v\n", err)

}
