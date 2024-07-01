package news

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type News interface {
	text(string)
	richText(string, string)
	Send(...string)
}

type FsNews struct {
	fsUrl string
}

// 检查接口实现
func _() {
	var _ News = (*FsNews)(nil)
}

func FsNew(fsUrl string) News {
	if fsUrl == "" {
		log.Printf("FsNew fsUrl empty, fsUrl:%v\n", fsUrl)
		return nil
	}
	return &FsNews{
		fsUrl: fsUrl,
	}
}

type fsInfo struct {
	MsgType string         `json:"msg_type"`
	Content map[string]any `json:"content"`
}

func (f *FsNews) Send(args ...string) {
	switch len(args) {
	case 0:
		log.Printf("rich text args length zero, args:%v\n", args)
		return
	case 1:
		f.text(args[0])
		return
	case 2:
		f.richText(args[0], args[1])
		return
	default:
		f.richText(args[0], args[1])
		return
	}
}

func (f *FsNews) text(msg string) {
	if msg == "" {
		log.Printf("FsNews text msg empty, msg:%v\n", msg)
		return
	}
	info := &fsInfo{
		MsgType: "text",
		Content: map[string]any{
			"text": msg,
		},
	}
	sendFsNews(f.fsUrl, info)
	return
}

func (f *FsNews) richText(title, msg string) {
	if msg == "" {
		log.Printf("FsNews rich text msg empty, msg:%v\n", msg)
		return
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
	return
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
