package news

// 自建应用 AccessToken 结构体
type InternalApp struct {
	AppAccessToken    string `json:"app_access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
}

// 自建应用 AccessToken 响应结构体
type InternalAccessTokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	AppAccessToken    string `json:"app_access_token"`
	Expire            int    `json:"expire"`
}

// 发送图片消息请求结构体
type ImageMessageReq struct {
	ReceiveID string `json:"receive_id"`
	MsgType   string `json:"msg_type"`
	Content   string `json:"content"`
	UUID      string `json:"uuid"`
}

type ImageInfo struct {
	ImageKey string `json:"image_key"`
}

// 发送图片消息响应结构体
type ImageMessageResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		MessageID  string `json:"message_id"`
		MsgType    string `json:"msg_type"`
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
		Deleted    bool   `json:"deleted"`
		Updated    bool   `json:"updated"`
		ChatID     string `json:"chat_id"`
		Sender     struct {
			ID         string `json:"id"`
			IDType     string `json:"id_type"`
			SenderType string `json:"sender_type"`
			TenantKey  string `json:"tenant_key"`
		} `json:"sender"`
		Body struct {
			Content string `json:"content"`
		} `json:"body"`
	} `json:"data"`
}