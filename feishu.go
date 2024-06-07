package news

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type News interface {
	Text(msg string)
	RichText(title, msg string)
}

type FsNews struct {
	fsUrl string
}

func _() {
	var _ News = (*FsNews)(nil)
}

func FsNew(fsUrl string) News {
	return &FsNews{
		fsUrl: fsUrl,
	}
}

type fsInfo struct {
	MsgType string         `json:"msg_type"`
	Content map[string]any `json:"content"`
}

func (f *FsNews) Text(msg string) {
	contentType := "application/json"
	info := &fsInfo{
		MsgType: "text",
		Content: map[string]any{
			"text": msg,
		},
	}

	data, _ := json.Marshal(info)

	resp, err := http.Post(f.fsUrl, contentType, strings.NewReader(string(data)))
	if err != nil {
		log.Printf("post failed, err:%v\n", err)
		return
	}

	_ = resp.Body.Close()
}

func (f *FsNews) RichText(title, msg string) {
	info := &fsInfo{
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
							// {
							// 	"tag":  "a",
							// 	"text": "请查看",
							// 	"href": "http://www.example.com/",
							// },
							// {
							// 	"tag":     "at",
							// 	"user_id": "all",
							// },
						},
					},
				},
			},
		},
	}
	data, _ := json.Marshal(info)

	resp, err := http.Post(f.fsUrl, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Printf("post failed, err:%v\n", err)
		return
	}

	_ = resp.Body.Close()
}
