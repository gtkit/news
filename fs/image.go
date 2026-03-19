package fs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gtkit/json"
)

// ImageKey 返回上传后的图片 key.
func (resp *UploadImageResp) ImageKey() string {
	return resp.Data.ImageKey
}

// UploadImage 上传图片到飞书，返回包含 image_key 的响应.
func (a *InternalApp) UploadImage(ctx context.Context, path string) (*UploadImageResp, error) {
	const api = "https://open.feishu.cn/open-apis/im/v1/images"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fs: open file %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	part, err := writer.CreateFormFile("image", filepath.Base(path))
	if err != nil {
		return nil, fmt.Errorf("fs: create form file: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("fs: copy file content: %w", err)
	}

	_ = writer.WriteField("image_type", "message")
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("fs: close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api, payload)
	if err != nil {
		return nil, fmt.Errorf("fs: create upload request: %w", err)
	}

	req.Header.Set("Authorization", a.authTenantToken())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fs: upload request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	const maxBody = 1 << 20
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return nil, fmt.Errorf("fs: read upload response: %w", err)
	}

	var uploadResp UploadImageResp
	if err := json.Unmarshal(data, &uploadResp); err != nil {
		return nil, fmt.Errorf("fs: decode upload response: %w", err)
	}

	if uploadResp.Code != 0 {
		return nil, fmt.Errorf("fs: upload image: code=%s, msg=%s",
			strconv.Itoa(uploadResp.Code), uploadResp.Msg)
	}

	return &uploadResp, nil
}

// DownloadImage 从飞书下载图片到指定路径.
func (a *InternalApp) DownloadImage(ctx context.Context, imageKey, path string) error {
	api := "https://open.feishu.cn/open-apis/im/v1/images/" + imageKey

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		return fmt.Errorf("fs: create download request: %w", err)
	}
	req.Header.Set("Authorization", a.authTenantToken())

	resp, err := defaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fs: download request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	const maxImage = 10 << 20
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxImage))
	if err != nil {
		return fmt.Errorf("fs: read image data: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("fs: write image file: %w", err)
	}

	return nil
}
