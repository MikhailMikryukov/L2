package parser

import (
	"L2.16/entities"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

// ParsePages получаем информацию о ресурсах
func (p *HTMLParser) ParsePages(pageURL *url.URL, depth int) ([]entities.Page, error) {
	// При глубине рекурсии 0 скачиваем только эту страницу
	if depth < 0 {
		return nil, fmt.Errorf("maximum recursion depth")
	}

	// Проверяем посещали ли уже
	if _, ok := p.visited[pageURL.String()]; ok {
		return nil, fmt.Errorf("alredy visited, %v", pageURL)
	}

	// Делаем ссылку абсолютной
	if !pageURL.IsAbs() {
		pageURL = p.BaseURL.ResolveReference(pageURL)
	}

	// Проверка принадлежит ли этому сайту
	if pageURL.Host != p.BaseURL.Host {
		return nil, fmt.Errorf("external domain: %s", pageURL.Host)
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", pageURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request error: %s", err)
	}

	// Выполняем запрос
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("new document error: %s", err)
	}

	// Помечаем как посещенный
	p.visited[pageURL.String()] = true

	// Получаем img, css, js файлы со страницы
	var pages []entities.Page
	page := p.getPageWithResources(doc)
	page.BaseHTMLURL = pageURL

	pages = append(pages, *page)

	// Проходимся по ссылкам и рекурсивно парсим страницы
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		if exists {
			linkURL, err := url.Parse(link)
			if err != nil {
				return
			}

			page, err := p.ParsePages(linkURL, depth-1)
			if err != nil {
				return
			}

			if page != nil {
				pages = append(pages, page...)
			}
		}
	})

	return pages, nil
}

func (p *HTMLParser) getPageWithResources(doc *goquery.Document) *entities.Page {
	page := entities.NewPage()

	// Получаем картинки
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		image, exists := selection.Attr("src")
		if exists {
			page.ImagesURL[image] = true
		}
	})

	// Получаем JS файлы
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		script, exists := selection.Attr("src")
		if exists {
			page.JsURL[script] = true

		}
	})

	// Получаем CSS файлы
	doc.Find("link").Each(func(i int, selection *goquery.Selection) {
		rel, relExists := selection.Attr("rel")
		if relExists {
			if rel == "stylesheet" {
				css, hrefExists := selection.Attr("href")
				if hrefExists {
					page.CssURL[css] = true
				}
			}
		}
	})

	mapSize := len(page.CssURL) + len(page.JsURL) + len(page.ImagesURL) + 1
	page.LocalPaths = make(map[string]string, mapSize)

	return page
}
