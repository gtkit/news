package news

import (
	"bytes"
	"context"
	"net/http"
	"strconv"

	"github.com/gtkit/json"

	"github.com/pkg/errors"
)

func (a *InternalApp) SendImageMessage(ctx context.Context, openID, imageKey string) error {
	api := "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=open_id"
	img, _ := json.Marshal(ImageInfo{
		ImageKey: imageKey,
	})
	msg := ImageMessageReq{
		ReceiveID: openID,
		MsgType:   "image",
		Content:   string(img),
	}
	payload, _ := json.Marshal(msg)

	// Send request to Feishu API
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", a.authTenantToken())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var respInfo ImageMessageResp
	if err := json.NewDecoder(resp.Body).Decode(&respInfo); err != nil {
		return err
	}

	// fmt.Printf("Send image message response: %+v", respInfo)
	if respInfo.Code != 0 {
		return errors.WithMessage(errors.New("错误码: "+strconv.Itoa(respInfo.Code)), respInfo.Msg)
	}

	return nil
}
