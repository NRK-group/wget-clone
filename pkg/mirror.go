package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func DownloadFile(URL string, rateLimit int64) ([]byte, error) {
	// Send GET request to the provided URL
	response, err := http.Get(URL)
	if err != nil {
		// Return nil and error if request fails
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		// Return nil and error if response status is not OK
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer

	// Read from the response body without rate limiting
	_, err = io.Copy(&data, response.Body)
	if err != nil {
		// Return nil and error if an error occurred
		return nil, err
	}

	// Return the bytes and a nil error if successful
	return data.Bytes(), nil
}

func FindStylesheet(doc *goquery.Document, url, folderName string) (rdoc *goquery.Document) {
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("href")
		if !exists {
			return
		} else {
			if strings.Contains(imgSrc, ".css") && strings.Contains(imgSrc, "https://") {
				MakeAFolder("./" + folderName + "/css/")
				fileName := path.Base(imgSrc)

				resp, err := DownloadFile(imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}

				if fileName[len(fileName)-4:] != ".css" && strings.Contains(fileName, ".css") {
					fileName = strings.SplitAfter(fileName, ".css")[0]
				}
				
				SaveBytesToFile("./"+folderName+"/css", fileName, resp)

				s.SetAttr("href", "./css/"+fileName)
			} else if strings.Contains(imgSrc, ".css") {
				fileName := path.Base(imgSrc)

				MakeAFolder("./" + folderName + "/css")

				resp, err := DownloadFile("https://"+folderName+"/"+imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}

				if fileName[len(fileName)-4:] != ".css" && strings.Contains(fileName, ".css") {
					fileName = strings.SplitAfter(fileName, ".css")[0]
				}

				SaveBytesToFile("./"+folderName+"/css", fileName, resp)

				s.SetAttr("href", "./css/"+fileName)
			}

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
				MakeAFolder("./" + folderName + "/js")
				fileName := path.Base(imgSrc)
				resp, err := DownloadFile(imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}
				SaveBytesToFile("./"+folderName+"/js", fileName, resp)

				s.SetAttr("src", "./js/"+fileName)
			} else if strings.Contains(imgSrc, ".js") {
				fileName := path.Base(imgSrc)
				MakeAFolder("./" + folderName + "/js")

				resp, err := DownloadFile("https://"+folderName+"/"+imgSrc, 0)
				if err != nil {
					fmt.Println(err)
					return
				}
				SaveBytesToFile("./"+folderName+"/js", fileName, resp)

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
						MakeAFolder("./" + folderName + "/img")
						fileName := path.Base(imgSrc)

						resp, err := DownloadFile(imgSrc, 0)
						if err != nil {
							fmt.Println(err)
							return
						}

						SaveBytesToFile("./"+folderName+"/img", fileName, resp)

						s.SetAttr("src", "./img/"+fileName)
					} else if strings.Contains(imgSrc, n) {
						fileName := path.Base(imgSrc)
						MakeAFolder("./" + folderName + "/img")

						resp, err := DownloadFile("https://"+folderName+"/"+imgSrc, 0)
						if err != nil {
							fmt.Println(err)
							return
						}
						SaveBytesToFile("./"+folderName+"/img", fileName, resp)
						s.SetAttr("src", "./img/"+fileName)
					}
				}
			}
		}

	})

	return doc
}

func FindUrlInStyle(doc *goquery.Document, url, folderName, Reject string) (rdoc string) {
	docString, _ := doc.Html()
	lines := strings.Split(docString, "\n")
	var returnStringArr []string
	var returnString []string
	var tagName string

	//
	for _, line := range lines {
		if !strings.Contains(line, "*") && strings.Contains(line, "background-image") {
			urls := strings.Split(strings.Split(line, ":")[1], ",")
			tagName = strings.Split(line, ":")[0] + ": "
			for _, url := range urls {
				if strings.Contains(url, "url('") {

					MakeAFolder("./" + folderName + "/img")

					resp, err := DownloadFile("http://"+folderName+"/"+strings.Split(url, `'`)[1][1:], 0)
					if err != nil {
						fmt.Println(err)
						return
					}
					SaveBytesToFile("./"+folderName+"/img", strings.Split(url, `'`)[1][1:], resp)

					returnString = append(returnString, "url("+`'`+"./img"+strings.Split(url, `'`)[1]+`')`)
				}
			}

			returnStringArr = append(returnStringArr, tagName+strings.Join(returnString, ",")+";")
		} else {
			returnStringArr = append(returnStringArr, line)
		}
	}

	return strings.Join(returnStringArr, "\n")
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
	res, err := http.Get(url) 
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
	if Exclude != "/css" {
		holder = FindStylesheet(doc, url, folderName)
	}
	if Exclude != "/js" {
		holder = Findjs(holder, url, folderName)
	}
	if Exclude != "/img" {
		holder = Findimg(holder, url, folderName, Reject)
		returnStr := FindUrlInStyle(holder, url, folderName, Reject)
		SaveBytesToFile("./"+folderName, "index.html", []byte(returnStr))
	} else {

		r, _ := holder.Html()
		SaveBytesToFile("./"+folderName, "index.html", []byte(r))
	}
}
