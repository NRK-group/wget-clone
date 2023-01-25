package pkg

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FindStylesheet(doc *goquery.Document) {
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("href")
		if !exists {
			fmt.Println("link href not found")
			return
		} else {
			if strings.Contains(imgSrc, ".css") && strings.Contains(imgSrc, "https://") {
				MakeAFolder("css")
				s.SetAttr("href", "css/"+"fileName")
			} else if strings.Contains(imgSrc, ".css") {
				MakeAFolder("css")
				s.SetAttr("href", "css/"+"fileName")
			}
			fmt.Println(imgSrc)
			// fmt.Println(doc.Html())
		}
	})
}

func Findjs(doc *goquery.Document) {
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if !exists {
			fmt.Println("script src not found")
			return
		} else {
			if strings.Contains(imgSrc, ".js") && strings.Contains(imgSrc, "https://") {
				MakeAFolder("js")
				s.SetAttr("src", "js/"+"fileName")
			} else if strings.Contains(imgSrc, ".js") {
				MakeAFolder("js")
				s.SetAttr("src", "js/"+"fileName")
			}
			fmt.Println(imgSrc)
			// fmt.Println(doc.Html())
		}
	})
}

func Findimg(doc *goquery.Document) {
	listImgSuffixes := []string{"jpg", "gif", "webb", "jpeg", "png"}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if !exists {
			fmt.Println("Image src not found")
			return
		} else {
			for _, n := range listImgSuffixes {
				if strings.Contains(imgSrc, n) && strings.Contains(imgSrc, "https://") {
					MakeAFolder("img")
					s.SetAttr("src", "img/"+"fileName")
				} else if strings.Contains(imgSrc, n) {
					MakeAFolder("img")
					s.SetAttr("src", "img/"+"fileName")
				}
			}
		}
	})
}

func MakeAFolder(name string) {
	if _, err := os.Stat(name); os.IsNotExist(err) { // check if a folder exist
		err := os.MkdirAll(name, 0o777)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
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

	FindStylesheet(doc)
	Findjs(doc)
	Findimg(doc)
}
