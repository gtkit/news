package news

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gtkit/json"
	"github.com/pkg/errors"
)

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
		body, _ := io.ReadAll(resp.Body)
		log.Print("Error in UploadImage Request: ", err)
		log.Print("Error in UploadImage Response: ", string(body))
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var uploadResp UploadImageResp
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		log.Print("Error in UploadImage Response: ", err)
		return nil, err
	}

	if uploadResp.Code != 0 {
		log.Print("Error in UploadImage Response: ", uploadResp.Msg)
		return nil, errors.New(uploadResp.Msg)
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
