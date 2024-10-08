package news

import (
	"bytes"
	"context"
	"net/http"

	"github.com/gtkit/json"
)

// InternalApp 内部应用, 用于获取内部应用的 access token
func NewInternalApp(ctx context.Context, appID, appSecret string) (*InternalApp, error) {
	api := "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"

	app := struct {
		AppID     string `json:"app_id"`
		AppSecret string `json:"app_secret"`
	}{
		AppID:     appID,
		AppSecret: appSecret,
	}

	// 发送请求
	body, _ := json.Marshal(app)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result InternalAccessTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	// fmt.Printf("InternalApp: %+v\n", result)

	return &InternalApp{
		AppAccessToken:    result.AppAccessToken,
		TenantAccessToken: result.TenantAccessToken,
		Expire:            result.Expire,
	}, nil
}

func (a *InternalApp) authTenantToken() string {
	return "Bearer " + a.TenantAccessToken
}

func (a *InternalApp) authAppToken() string {
	return "Bearer " + a.AppAccessToken
}
