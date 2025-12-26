package test_images

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Downloader struct {
	client *http.Client
}

type Option func(*Downloader)

func WithTimeout(d time.Duration) Option {
	return func(dl *Downloader) {
		dl.client.Timeout = d
	}
}

func New(opts ...Option) *Downloader {
	d := &Downloader{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (d *Downloader) Download(
	ctx context.Context,
	url string,
	filePath string,
) (int64, error) {

	if url == "" {
		return 0, ErrEmptyURL
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("http error: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, resp.Body)
}
