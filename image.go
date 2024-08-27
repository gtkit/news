package news

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gtkit/json"
	"github.com/pkg/errors"
)

type UploadImageResp struct {
	Code int `json:"code"`
	Data struct {
		ImageKey string `json:"image_key"`
	} `json:"data"`
	Msg   string `json:"msg"`
	Error any    `json:"error"`
	Helps any    `json:"helps"`
}

func (a *InternalApp) UploadImage(ctx context.Context, path string) (*UploadImageResp, error) {
	api := "https://open.feishu.cn/open-apis/im/v1/images"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	part, errFile := writer.CreateFormFile("image", filepath.Base(path))

	_, errFile = io.Copy(part, file)
	if errFile != nil {
		return nil, errFile
	}

	_ = writer.WriteField("image_type", "message")
	if err = writer.Close(); err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api, payload)

	// r.Header.Add("Content-Type", "multipart/form-data; boundary=---7MA4YWxkTrZu0gW")
	req.Header.Add("Authorization", a.authTenantToken())
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var uploadResp UploadImageResp
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	if uploadResp.Code != 0 {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}
	return &uploadResp, nil
}

func (a *InternalApp) DownloadImage(ctx context.Context, imageKey, path string) error {
	api := "https://open.feishu.cn/open-apis/im/v1/images/" + imageKey

	client := &http.Client{}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	req.Header.Add("Authorization", a.authTenantToken())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	resbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, resbytes, 0644); err != nil {
		return err
	}

	return nil
}
