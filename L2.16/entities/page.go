package entities

import "net/url"

// Page структура одной страницы сайта
type Page struct {
	BaseHTMLURL *url.URL
	ImagesURL   map[string]bool
	CssURL      map[string]bool
	JsURL       map[string]bool
	LocalPaths  map[string]string
}

// NewPage инициализируем мапы хранения ссылок
func NewPage() *Page {
	return &Page{
		ImagesURL: make(map[string]bool),
		CssURL:    make(map[string]bool),
		JsURL:     make(map[string]bool),
	}
}
