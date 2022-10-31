// @Author xiaozhaofu 2022/10/31 20:10:00
package news

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func FsWarnText(fsurl, msg string) {

	contentType := "application/json"

	info := &fsInfo{
		MsgType: "text",
		Content: map[string]string{
			"text": msg,
		},
	}

	data, _ := json.Marshal(info)

	resp, err := http.Post(fsurl, contentType, strings.NewReader(string(data)))
	if err != nil {
		log.Printf("post failed, err:%v\n", err)
		return
	}

	resp.Body.Close()

}

type fsInfo struct {
	MsgType string            `json:"msg_type"`
	Content map[string]string `json:"content"`
}
