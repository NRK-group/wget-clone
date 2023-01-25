package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)



func FindStylesheet( doc *goquery.Document) {
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("href")
		if !exists {
			fmt.Println("link href not found")
			return
		} else {
			if strings.Contains(imgSrc, ".css") &&  strings.Contains(imgSrc, "https://"){

				//s.SetAttr("href", "keivon") 
			} else if strings.Contains(imgSrc, ".css") {

			}
			fmt.Println(imgSrc)
			//fmt.Println(doc.Html())
		}
	})
}

func Findjs( doc *goquery.Document) {

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if !exists {
			fmt.Println("script src not found")
			return
		} else {
			if strings.Contains(imgSrc, ".js") &&  strings.Contains(imgSrc, "https://"){

				//s.SetAttr("href", "keivon") 
			} else if strings.Contains(imgSrc, ".js") {

			}
			fmt.Println(imgSrc)
			//fmt.Println(doc.Html())
		}
	})

}


func Findimg( doc *goquery.Document) {

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if !exists {
			fmt.Println("Image src not found")
			return
		} else {
			if strings.Contains(imgSrc, ".css") &&  strings.Contains(imgSrc, "https://"){

				//s.SetAttr("href", "keivon") 
			} else if strings.Contains(imgSrc, ".css") {

			}
			fmt.Println(imgSrc)
			//fmt.Println(doc.Html())
		}
	})
}


func mirror(url string) {
	res, err := http.Get("https://jonathanmh.com/web-scraping-golang-goquery/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	FindStylesheet(doc);
	Findjs(doc)
	Findimg(doc)

	
}
