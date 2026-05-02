package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// DownloadResource метод для скачивания ресурсов, возвращает путь до сохраненного файла
func (d *Downloader) DownloadResource(resourceURL string, ext string) (string, error) {

	parsed, err := url.Parse(resourceURL)
	if err != nil {
		return "", fmt.Errorf("parse error: %s", err)
	}

	// Делаем абсолютный URL если нужно
	if !parsed.IsAbs() {
		parsed = d.BaseURL.ResolveReference(parsed)
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", parsed.String(), nil)
	if err != nil {
		return "", fmt.Errorf("creating request error: %v", err)
	}

	// Выполняем запрос
	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	// Создаем путь к файлу
	relativePath := getRelativePath(*parsed, ext)
	savePath := filepath.Join(d.BaseDir, relativePath)

	// Создаем директорию если нужно
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return "", fmt.Errorf("creating directions error: %v", err)
	}

	// Создаем файл
	file, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("creating file error: %v", err)
	}
	defer file.Close()

	// Копируем данные
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("copy data error: %v", err)
	}

	return relativePath, nil
}

func getRelativePath(parsedURL url.URL, ext string) string {
	var relative string

	switch ext {
	case "html":
		if parsedURL.Path == "" || parsedURL.Path == "/" {
			parsedURL.Path = "index.html"
		}
		relative = filepath.Join("html", parsedURL.Path)
	case "css":
		relative = filepath.Join("css", filepath.Base(parsedURL.Path))
	case "img":
		relative = filepath.Join("img", filepath.Base(parsedURL.Path))
	case "js":
		relative = filepath.Join("js", filepath.Base(parsedURL.Path))
	}

	return relative
}
