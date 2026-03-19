package fs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gtkit/json"
)

const (
	// FSMsgAPI 飞书发送消息 API 地址.
	FSMsgAPI = "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=open_id"
)

// SendImageMsg 发送图片消息.
func (a *InternalApp) SendImageMsg(ctx context.Context, openID, filepath string) error {
	upimg, err := a.UploadImage(ctx, filepath)
	if err != nil {
		return fmt.Errorf("fs: upload image: %w", err)
	}

	img, err := json.Marshal(ImageInfo{ImageKey: upimg.ImageKey()})
	if err != nil {
		return fmt.Errorf("fs: marshal image info: %w", err)
	}

	payload, err := json.Marshal(MessageReq{
		ReceiveID: openID,
		MsgType:   "image",
		Content:   string(img),
	})
	if err != nil {
		return fmt.Errorf("fs: marshal message: %w", err)
	}

	return doRequest(ctx, a.authTenantToken(), payload)
}

// SendTextMsg 发送文本消息.
func (a *InternalApp) SendTextMsg(ctx context.Context, openID, msg string) error {
	tm, err := json.Marshal(TextInfo{Text: msg})
	if err != nil {
		return fmt.Errorf("fs: marshal text info: %w", err)
	}

	payload, err := json.Marshal(MessageReq{
		ReceiveID: openID,
		MsgType:   "text",
		Content:   string(tm),
	})
	if err != nil {
		return fmt.Errorf("fs: marshal message: %w", err)
	}

	return doRequest(ctx, a.authTenantToken(), payload)
}

func doRequest(ctx context.Context, token string, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, FSMsgAPI, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("fs: create request: %w", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fs: send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	const maxBody = 1 << 20
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return fmt.Errorf("fs: read response: %w", err)
	}

	var respInfo MessageResp
	if err := json.Unmarshal(data, &respInfo); err != nil {
		return fmt.Errorf("fs: decode response: %w", err)
	}

	if respInfo.Code != 0 {
		return fmt.Errorf("fs: api error: code=%s, msg=%s", strconv.Itoa(respInfo.Code), respInfo.Msg)
	}

	return nil
}
