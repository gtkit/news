package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type fsInfo struct {
	MsgType string         `json:"msg_type"`
	Content map[string]any `json:"content"`
}

// WebHookSend 发送飞书机器人 webhook 消息.
// 1. 发送文本消息：WebHookSend(url, "hello world").
// 2. 发送富文本消息：WebHookSend(url, "标题", "正文", "提示", "超链接").
func WebHookSend(url string, args ...string) error {
	return WebHookSendCtx(context.Background(), url, args...)
}

// WebHookSendCtx 发送飞书机器人 webhook 消息（带 context）.
func WebHookSendCtx(ctx context.Context, url string, args ...string) error {
	if url == "" {
		return fmt.Errorf("fs: webhook url is empty")
	}

	var info *fsInfo
	switch len(args) {
	case 0:
		return fmt.Errorf("fs: webhook args is empty")
	case 1:
		info = textMsg(args[0])
	case 2:
		info = richTextMsg(args[0], args[1], "", "")
	default:
		tips, href := "", ""
		if len(args) >= 3 {
			tips = args[2]
		}
		if len(args) >= 4 {
			href = args[3]
		}
		info = richTextMsg(args[0], args[1], tips, href)
	}
	if info == nil {
		return fmt.Errorf("fs: build message failed")
	}

	return sendFsNews(ctx, url, info)
}

func textMsg(msg string) *fsInfo {
	if msg == "" {
		return nil
	}
	return &fsInfo{
		MsgType: "text",
		Content: map[string]any{"text": msg},
	}
}

func richTextMsg(title, msg, tips, hyperlink string) *fsInfo {
	if msg == "" {
		return nil
	}

	elements := []map[string]any{
		{"tag": "text", "text": msg},
	}
	if hyperlink != "" && tips != "" {
		elements = append(elements, map[string]any{
			"tag": "a", "text": tips, "href": hyperlink,
		})
	}
	elements = append(elements, map[string]any{
		"tag": "at", "user_id": "all",
	})

	return &fsInfo{
		MsgType: "post",
		Content: map[string]any{
			"post": map[string]any{
				"zh_cn": map[string]any{
					"title":   title,
					"content": []any{elements},
				},
			},
		},
	}
}

func sendFsNews(ctx context.Context, fsURL string, info *fsInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("fs: marshal webhook message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fsURL, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("fs: create webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fs: webhook request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	const maxBody = 1 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return fmt.Errorf("fs: read webhook response: %w", err)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("fs: decode webhook response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("fs: webhook api error: code=%d, msg=%s", result.Code, result.Msg)
	}

	return nil
}
