package downloader

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// Downloader загрузчик страниц
type Downloader struct {
	BaseURL    *url.URL
	BaseDir    string
	client     *http.Client
	Downloaded map[string]bool
}

// NewDownloader устанавливаем папку для сохранения, базовую ссылку и клиент
func NewDownloader(baseURL *url.URL, baseDir string) (*Downloader, error) {
	savePath := path.Join(baseDir, baseURL.Host)
	err := os.MkdirAll(savePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("directory create err: %s", err)
	}

	return &Downloader{
		BaseURL: baseURL,
		BaseDir: savePath,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    100,
				IdleConnTimeout: 30 * time.Second,
			},
		},
	}, nil
}
