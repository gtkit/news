package fs

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type fsInfo struct {
	MsgType string         `json:"msg_type"`
	Content map[string]any `json:"content"`
}

// 发送飞书机器人webhook消息.
// 1. 发送文本消息：WebHookSend(url, "hello world")
// 2. 发送富文本消息：WebHookSend(url, "标题", "正文", "提示", "超链接")
func WebHookSend(url string, args ...string) {
	var info *fsInfo
	switch len(args) {
	case 0:
		log.Printf("rich text args length zero, args:%v\n", args)
		return
	case 1:
		info = text(args[0])
	case 2:
		info = richText(args[0], args[1], "", "")
	case 4:
		info = richText(args[0], args[1], args[2], args[3])
	default:
		info = richText(args[0], args[1], args[2], args[3])

	}
	if info == nil {
		return
	}
	sendFsNews(url, info)
}

// text：普通文本.
func text(msg string) *fsInfo {
	if msg == "" {
		log.Printf("FsNews text msg empty, msg:%v\n", msg)
		return nil
	}
	return &fsInfo{
		MsgType: "text",
		Content: map[string]any{
			"text": msg,
		},
	}
}

// text：普通文本.
// a：超链接.
// at：@符号.
// img：图.
func richText(title, msg, tips, hyperlink string) *fsInfo {
	if msg == "" {
		log.Printf("FsNews rich text msg empty, msg:%v\n", msg)
		return nil
	}
	if hyperlink == "" || tips == "" {
		return &fsInfo{
			MsgType: "post",
			Content: map[string]any{
				"post": map[string]any{
					"zh_cn": map[string]any{
						"title": title,
						"content": []any{
							[]map[string]any{
								{
									"tag":  "text",
									"text": msg,
								},
								{
									"tag":     "at",
									"user_id": "all",
								},
							},
						},
					},
				},
			},
		}
	}
	return &fsInfo{
		MsgType: "post",
		Content: map[string]any{
			"post": map[string]any{
				"zh_cn": map[string]any{
					"title": title,
					"content": []any{
						[]map[string]any{
							{
								"tag":  "text",
								"text": msg,
							},
							{
								"tag":  "a",
								"text": tips,
								"href": hyperlink,
							},
							{
								"tag":     "at",
								"user_id": "all",
							},
						},
					},
				},
			},
		},
	}
}

func sendFsNews(fsUrl string, info *fsInfo) {
	data, err := json.Marshal(info)
	if err != nil {
		log.Printf("sendFsNews marshal failed, err:%v\n", err)
		return
	}

	resp, err := http.Post(fsUrl, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Printf("sendFsNews post failed, err:%v\n", err)
		return
	}

	_ = resp.Body.Close()
}
