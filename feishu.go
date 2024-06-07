package news

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type News interface {
	Text(msg string)
	RichText(...string)
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
	info := &fsInfo{
		MsgType: "text",
		Content: map[string]any{
			"text": msg,
		},
	}
	sendFsNews(f.fsUrl, info)
	return
}

func (f *FsNews) RichText(args ...string) {
	var title, msg string
	switch len(args) {
	case 0:
		log.Printf("rich text args length zero, args:%v\n", args)
		return
	case 1:
		msg = args[0]
	case 2:
		title, msg = args[0], args[1]
	default:
		title, msg = args[0], args[1]
	}

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
	sendFsNews(f.fsUrl, info)
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
