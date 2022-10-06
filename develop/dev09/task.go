package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Flags struct {
	lDepth int
	url    string
	path   string
}

// скачать страницу
func getUrlPage(uri string) []byte {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Fprintln(os.Stderr, "http err: ", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error status: %s\n", resp.Status)
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading page: %v\n", err)
		return nil
	}
	return b
}

// запуск wget
func wget(f Flags) {
	fmt.Println("downloading " + f.url)
	if f.lDepth < 1 {
		fmt.Fprintln(os.Stderr, "wrong depth")
		return
	}
	wgetRec(f, f.url, "", 0)
}

// рекурсивный wget
func wgetRec(f Flags, url string, oldPath string, recDepth int) {
	if recDepth == f.lDepth {
		return
	}
	page := getUrlPage(url)
	if page != nil {
		links := parseLinks(page)
		//находим ссылки, проходим по ним и меняем на соответствующие директории
		for _, val := range links {
			byteVal := []byte(val)
			newLink := "./" + oldPath + linkToFilePath(val)
			page = bytes.ReplaceAll(page, byteVal, []byte(newLink))
		}
		//записываем страницу
		err := writeToFile(page, oldPath+`/`+linkToFilePath(url)+"/")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error writing to file")
			return
		}
		for _, v := range links {
			link := `https://` + linkToFilePath(v)
			wgetRec(f, link, linkToFilePath(url)+`/`, recDepth+1)
		}
	}
}

// запись в файл данных
func writeToFile(data []byte, dirname string) error {
	os.Mkdir("./"+dirname+"/", 0666)
	f, err := os.Create("./" + dirname + "/" + "index.html")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// возвращает слайс ссылок со страницы
func parseLinks(data []byte) []string {
	links := make([]string, 0)
	body := bytes.NewReader(data)
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
	return links
}

// парсим ссылку, возвращаем путь к фалй без слешей
func linkToFilePath(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		fmt.Fprintln(os.Stderr, "problem while parsing url")
		return ""
	}
	p := u.Hostname() + "/" + u.EscapedPath()
	fullPath := strings.Split(p, "/")
	if len(fullPath[len(fullPath)-1]) == 0 {
		fullPath = fullPath[:len(fullPath)-1]
	}
	return path.Join(fullPath...)
}
func flagsInit() *Flags {
	f := &Flags{}
	flag.StringVar(&f.url, "f", "\t", "url")
	flag.IntVar(&f.lDepth, "l", 1, "depth of download")
	flag.Parse()
	return f
}

func main() {
	//пример запуска go run . -f 'https://wikipedia.org' -l 2
	wget(*flagsInit())
}
