package pkg

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FindStylesheet(doc *goquery.Document, url, folderName string) (rdoc *goquery.Document) {
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("href")
		if !exists {
			return
		} else {
			if strings.Contains(imgSrc, ".css") && strings.Contains(imgSrc, "https://") {
				MakeAFolder("./"+folderName +"/css/")
				fileName := path.Base(imgSrc)

				resp, err := DownloadFile(imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}
				SaveBytesToFile("./"+folderName +"/css/"+fileName, resp)

				s.SetAttr("href", "./css/"+fileName)
			} else if strings.Contains(imgSrc, ".css") {
				fileName := path.Base(imgSrc)
				MakeAFolder("./"+folderName +"/css")

				resp, err := DownloadFile("https://"+folderName+imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}
				if(fileName[len(fileName)-4:] != ".css"){
					fileName = fileName+".css"
				}

				SaveBytesToFile("./"+folderName +"/css/"+fileName, resp)
				
				s.SetAttr("href", "./css/"+fileName)
			}
			// fmt.Println(imgSrc)
		}
	})

	return doc
}

func Findjs(doc *goquery.Document, url, folderName string) (rdoc *goquery.Document) {
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if !exists {
			return
		} else {
			if strings.Contains(imgSrc, ".js") && strings.Contains(imgSrc, "https://") {
				MakeAFolder("./"+folderName +"/js")
				fileName := path.Base(imgSrc)
				resp, err := DownloadFile(imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}
				SaveBytesToFile("./"+folderName +"/js/"+fileName, resp)

				s.SetAttr("src", "./js/"+fileName)
			} else if strings.Contains(imgSrc, ".js") {
				fileName := path.Base(imgSrc)
				MakeAFolder("./"+folderName +"/js")

				resp, err := DownloadFile(url+imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}
				SaveBytesToFile("./"+folderName +"/js/"+fileName, resp)
				
				s.SetAttr("src", "./js/"+fileName)
			}
			// fmt.Println(imgSrc)
		}
	})

	return doc
}

func Findimg(doc *goquery.Document, url, folderName, Reject string) (rdoc *goquery.Document) {
	listImgSuffixes := []string{"jpg", "gif", "webb", "jpeg", "png"}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if !exists {
			return
		} else {
			for _, n := range listImgSuffixes {
				if Reject != n {
				if strings.Contains(imgSrc, n) && strings.Contains(imgSrc, "https://") {
					MakeAFolder("./"+folderName +"/img")
					fileName := path.Base(imgSrc)

					resp, err := DownloadFile(imgSrc, 0)
					if err != nil {
						fmt.Println(err)
						return
					}

					SaveBytesToFile("./"+folderName +"/img/"+fileName, resp)

					s.SetAttr("src", "./img/"+fileName)
				} else if strings.Contains(imgSrc, n) {
					fileName := path.Base(imgSrc)
					MakeAFolder("./"+folderName +"/img")

					resp, err := DownloadFile(url+imgSrc, 0)
					if err != nil {
						fmt.Println(err)
						return
					}
					SaveBytesToFile("./"+folderName +"/img/"+fileName, resp)
					s.SetAttr("src", "./img/"+fileName)
				}
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

func Mirror(url, Exclude, Reject string) {
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

	folderName := strings.Split(url, "/")[2]
	MakeAFolder(folderName)
	var holder *goquery.Document
	if (Exclude != "/css") {
	holder = FindStylesheet(doc, url, folderName)
	}
	if (Exclude != "/js") {
	holder = Findjs(holder, url, folderName)
	}
	if (Exclude != "/img") {
	holder = Findimg(holder, url, folderName, Reject)
	}
	r, _ := holder.Html()
	// fmt.Println(holder.Html())
	SaveBytesToFile("./"+folderName +"/index.html", []byte(r))
}
