package parser

import (
	"net/http"
	"net/url"
	"time"
)

// HTMLParser парсер
type HTMLParser struct {
	BaseURL *url.URL
	visited map[string]bool
	Client  *http.Client
}

// NewParser создаем клиент и устанавливаем базовую ссылку откуда начинаем парсить
func NewParser(baseURL *url.URL) (*HTMLParser, error) {
	return &HTMLParser{
		BaseURL: baseURL,
		visited: make(map[string]bool),
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    100,
				IdleConnTimeout: 30 * time.Second,
			},
		},
	}, nil
}
