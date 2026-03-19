package fs

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/gtkit/json"
	"github.com/gtkit/news"
)

const (
	// AccessTokenAPI 飞书获取 tenant access token 的 API 地址.
	AccessTokenAPI = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
)

// defaultClient 复用的 HTTP 客户端，避免每次请求新建.
var defaultClient = &http.Client{}

// EmptyApp 返回一个空的 InternalApp 实例.
func EmptyApp() *InternalApp {
	return &InternalApp{}
}

// NewInternalApp 内部应用, 用于获取内部应用的 access token.
func NewInternalApp(ctx context.Context, appID, appSecret string, cache ...news.AppCacher) (*InternalApp, error) {
	app := struct {
		AppID     string `json:"app_id"`
		AppSecret string `json:"app_secret"`
	}{
		AppID:     appID,
		AppSecret: appSecret,
	}

	body, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("fs: marshal app request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, AccessTokenAPI, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("fs: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fs: send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result InternalAccessTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("fs: decode response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("fs: api error: code=%d, msg=%s", result.Code, result.Msg)
	}

	internalapp := &InternalApp{
		AppAccessToken:    result.AppAccessToken,
		TenantAccessToken: result.TenantAccessToken,
		Expire:            result.Expire,
	}

	if len(cache) > 0 && cache[0] != nil {
		cache[0].Set(internalapp)
	}

	return internalapp, nil
}

func (a *InternalApp) authTenantToken() string {
	return "Bearer " + a.TenantAccessToken
}

func (a *InternalApp) authAppToken() string {
	return "Bearer " + a.AppAccessToken
}

// Expired 返回 token 的过期时间（秒）.
func (a *InternalApp) Expired() int {
	return a.Expire
}
