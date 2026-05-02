package main

import (
	"L2.16/downloader"
	"L2.16/parser"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("format: ./wget [link] [recursion depth] [save path]")
		return
	}
	
	// Получаем ссылку на скачивание
	parsedUrl, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println("url parsing error", err)
		return
	}

	// Получаем глубину рекурсии
	depth, err := strconv.Atoi(os.Args[2])
	if err != nil {
		depth = 0
	}

	// Папка куда сохраняем сайт
	savePath := os.Args[3]

	// Создаем Html parser
	p, err := parser.NewParser(parsedUrl)
	if err != nil {
		fmt.Println("new html parser creation error", err)
		return
	}

	fmt.Println("parsing pages...")
	// Получаем срез страниц для скачивания
	pages, err := p.ParsePages(parsedUrl, depth)
	if err != nil {
		fmt.Println("parse pages error", err)
		return
	}

	// Создаем загрузчик
	d, err := downloader.NewDownloader(parsedUrl, savePath)
	if err != nil {
		fmt.Println("new downloader creation error", err)
		return
	}

	fmt.Println("downloading pages...")

	// Скачиваем страницы
	err = d.DownloadPages(pages)
	if err != nil {
		fmt.Println("download page error", err)
	}

	fmt.Println("changing links to local...")

	// Заменяем ссылки на локальные на ресурсы
	for _, page := range pages {
		err = page.ReplaceLinks(d.BaseDir, d.Downloaded)
		if err != nil {
			fmt.Println(err)
		}
	}
}
