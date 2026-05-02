package downloader

import (
	"L2.16/entities"
	"fmt"
)

// DownloadPages скачиваем весь срез страниц
func (d *Downloader) DownloadPages(pages []entities.Page) error {
	d.Downloaded = make(map[string]bool, len(pages))
	for _, page := range pages {
		link := page.BaseHTMLURL.String()
		_, downloaded := d.Downloaded[link]
		if !downloaded {
			fmt.Println("Downloading ", link)
			err := d.downloadPage(&page)
			if err != nil {
				return err
			}
			d.Downloaded[link] = true
		}
	}

	return nil
}

func (d *Downloader) downloadPage(page *entities.Page) error {

	// Сохраняем саму страницу
	savedRelativePath, err := d.DownloadResource(page.BaseHTMLURL.String(), "html")
	if err != nil {
		return fmt.Errorf("download html error: %v", err)
	}
	page.LocalPaths[page.BaseHTMLURL.String()] = savedRelativePath

	// Сохраняем css файлы
	for css := range page.CssURL {
		savedRelativePath, err = d.DownloadResource(css, "css")
		if err != nil {
			return fmt.Errorf("download css error: %v", err)
		}
		page.LocalPaths[css] = savedRelativePath
	}

	// Сохраняем картинки
	for img := range page.ImagesURL {
		savedRelativePath, err = d.DownloadResource(img, "img")
		if err != nil {
			return fmt.Errorf("download img error: %v", err)
		}
		page.LocalPaths[img] = savedRelativePath
	}

	// Сохраняем JS файлы
	for js := range page.JsURL {
		savedRelativePath, err = d.DownloadResource(js, "js")
		if err != nil {
			return fmt.Errorf("download js error: %v", err)
		}
		page.LocalPaths[js] = savedRelativePath
	}

	return nil
}
