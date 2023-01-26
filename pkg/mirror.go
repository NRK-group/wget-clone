package pkg

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FindStylesheet(doc *goquery.Document, url string) (rdoc *goquery.Document) {
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("href")
		if !exists {
			fmt.Println("link href not found")
			return
		} else {
			if strings.Contains(imgSrc, ".css") && strings.Contains(imgSrc, "https://") {
				MakeAFolder("css")
				fileName := path.Base(imgSrc)
				s.SetAttr("href", "css/"+fileName)
			} else if strings.Contains(imgSrc, ".css") {
				fileName := path.Base(imgSrc)
				MakeAFolder("css")
				s.SetAttr("href", "css/"+fileName)
			}
			// fmt.Println(imgSrc)

		}
	})

	return doc
}

func Findjs(doc *goquery.Document, url string) (rdoc *goquery.Document) {
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if !exists {
			fmt.Println("script src not found")
			return
		} else {
			if strings.Contains(imgSrc, ".js") && strings.Contains(imgSrc, "https://") {
				MakeAFolder("js")
				fileName := path.Base(imgSrc)
				s.SetAttr("src", "js/"+fileName)
			} else if strings.Contains(imgSrc, ".js") {
				fileName := path.Base(imgSrc)
				MakeAFolder("js")
				s.SetAttr("src", "js/"+fileName)
			}
			// fmt.Println(imgSrc)

		}
	})

	return doc
}

func Findimg(doc *goquery.Document, url string) (rdoc *goquery.Document) {
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
					fileName := path.Base(imgSrc)
					s.SetAttr("src", "img/"+fileName)
				} else if strings.Contains(imgSrc, n) {
					fileName := path.Base(imgSrc)
					MakeAFolder("img")
					s.SetAttr("src", "img/"+fileName)
				}
			}
		}
		// fmt.Println(imgSrc)
	})

	return doc
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
	res, err := http.Get(url) //"https://jonathanmh.com/web-scraping-golang-goquery/"
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

	holder := FindStylesheet(doc, url)
	holder = Findjs(holder, url)
	holder = Findimg(holder, url)
	holder.Html()
}
