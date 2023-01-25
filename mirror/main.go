package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	response, err := http.Get("https://github.com/PuerkitoBio/goquery")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	file, err := os.Create("example.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}

	// Parse the HTML and find the URLs of the CSS and JavaScript files
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("link[rel='stylesheet']").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			downloadFile(url)
		}
	})
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if ok {
			downloadFile(url)
		}
	})
}

func downloadFile(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	file, err := os.Create(getFileName(url))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err)
	}
}

func getFileName(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}
