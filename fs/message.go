package fs

import (
	"bytes"
	"context"
	"net/http"
	"strconv"

	"github.com/gtkit/json"

	"github.com/pkg/errors"
)

const (
	FSMsgAPI = "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=open_id"
)

// 发送图片消息
func (a *InternalApp) SendImageMsg(ctx context.Context, openID, filepath string) error {
	// 上传图片
	upimg, err := a.UploadImage(ctx, filepath)
	if err != nil {
		return errors.WithMessage(err, "飞书 upload image failed")
	}

	img, _ := json.Marshal(ImageInfo{
		ImageKey: upimg.ImageKey(),
	})
	msg := MessageReq{
		ReceiveID: openID,
		MsgType:   "image",
		Content:   string(img),
	}
	payload, _ := json.Marshal(msg)

	return doRequest(ctx, a.authTenantToken(), payload)
}

// 发送文本消息
func (a *InternalApp) SendTextMsg(ctx context.Context, openID, msg string) error {
	// todo: implement SendTextMessage
	tm, _ := json.Marshal(TextInfo{
		Text: msg,
	})
	reqinfo := MessageReq{
		ReceiveID: openID,
		MsgType:   "text",
		Content:   string(tm),
	}
	payload, _ := json.Marshal(reqinfo)

	return doRequest(ctx, a.authTenantToken(), payload)
}

func doRequest(ctx context.Context, token string, payload []byte) error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, FSMsgAPI, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var respInfo MessageResp
	if err := json.NewDecoder(resp.Body).Decode(&respInfo); err != nil {
		return err
	}

	// fmt.Printf("Send image message response: %+v", respInfo)
	if respInfo.Code != 0 {
		return errors.WithMessage(errors.New("错误码: "+strconv.Itoa(respInfo.Code)), respInfo.Msg)
	}

	return nil
}
