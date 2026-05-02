package entities

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ReplaceLinks меняем ссылки в документе на локальные
func (p *Page) ReplaceLinks(baseDir string, downloaded map[string]bool) error {

	relativeMainHTMLPath := p.LocalPaths[p.BaseHTMLURL.String()]
	absoluteMainHTMLPath := filepath.Join(baseDir, relativeMainHTMLPath)

	htmlFile, err := os.Open(absoluteMainHTMLPath)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(htmlFile))
	if err != nil {
		return err
	}

	// Заменяем ссылки на картинки
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		image, exists := selection.Attr("src")
		if exists {
			if localPath, ok := p.LocalPaths[image]; ok {
				// Получаем относительный путь от html страницы
				relPath, err := getRelativePathForResource(relativeMainHTMLPath, localPath)
				if err != nil {
					fmt.Println(err)
				}
				// Изменяем путь у атрибута
				selection.SetAttr("src", relPath)
			}
		}
	})

	// Заменяем ссылки на js файлы
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		script, exists := selection.Attr("src")
		if exists {
			if localPath, ok := p.LocalPaths[script]; ok {
				// Получаем относительный путь от html страницы
				relPath, err := getRelativePathForResource(relativeMainHTMLPath, localPath)
				if err != nil {
					fmt.Println(err)
				}
				// Изменяем путь у атрибута
				selection.SetAttr("src", relPath)
			}
		}
	})

	// Заменяем ссылки на css файлы
	doc.Find("link").Each(func(i int, selection *goquery.Selection) {
		rel, exists := selection.Attr("rel")
		if exists {
			if rel == "stylesheet" {
				css, exists := selection.Attr("href")
				if exists {
					if localPath, ok := p.LocalPaths[css]; ok {
						// Получаем относительный путь от html страницы
						relPath, err := getRelativePathForResource(relativeMainHTMLPath, localPath)
						if err != nil {
							fmt.Println(err)
						}
						// Изменяем путь у атрибута
						selection.SetAttr("href", relPath)
					}
				}
			}

		}
	})

	// Заменяем ссылки на другие html файлы
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		if exists {
			linkURL, err := url.Parse(link)
			if err != nil {
				return
			}
			// Делаем ссылку абсолютной
			if !linkURL.IsAbs() {
				linkURL = p.BaseHTMLURL.ResolveReference(linkURL)
			}

			// Меняем ссылки только если скачивали эти страницы
			_, ok := downloaded[linkURL.String()]

			if !ok {
				return
			}

			var relPath string

			if linkURL.Path == "/" {
				linkURL.Path = "index.html"
			}

			var linkFilePath string
			linkFilePath = "html" + string(filepath.Separator) + linkURL.Path
			if !strings.HasSuffix(linkURL.Path, "index.html") {
				linkFilePath += ".html"
			}
			// Получаем относительный путь от этой html страницы
			relPath, err = getRelativePathForResource(relativeMainHTMLPath, linkFilePath)
			if err != nil {
				fmt.Println(err)
			}

			selection.SetAttr("href", relPath)

		}
	})

	// Перезаписываем измененный html файл
	result, err := doc.Html()
	if err != nil {
		return err
	}

	err = os.WriteFile(absoluteMainHTMLPath, []byte(result), 0666)
	if err != nil {
		return err
	}

	return nil
}

func getRelativePathForResource(htmlPath string, resourcePath string) (string, error) {
	dirs := filepath.Dir(htmlPath)
	return filepath.Rel(dirs, resourcePath)
}
