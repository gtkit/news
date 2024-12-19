package fs

// 自建应用 AccessToken 结构体
type InternalApp struct {
	AppAccessToken    string `json:"app_access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

// 自建应用 AccessToken 响应结构体
type InternalAccessTokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	AppAccessToken    string `json:"app_access_token"`
	Expire            int    `json:"expire"`
}

// 发送消息请求结构体
type MessageReq struct {
	ReceiveID string `json:"receive_id"`
	MsgType   string `json:"msg_type"`
	Content   string `json:"content"`
	UUID      string `json:"uuid"`
}

type ImageInfo struct {
	ImageKey string `json:"image_key"`
}

type TextInfo struct {
	Text string `json:"text"`
}

// 发送消息响应结构体
type MessageResp struct {
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

type UploadImageResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ImageKey string `json:"image_key"`
	} `json:"data"`
	Error struct {
		Message              string `json:"message"`
		LogID                string `json:"log_id"`
		PermissionViolations []struct {
			Type    string `json:"type"`
			Subject string `json:"subject"`
		} `json:"permission_violations"`
		Helps []struct {
			URL         string `json:"url"`
			Description string `json:"description"`
		} `json:"helps"`
	} `json:"error,omitempty"`
}
