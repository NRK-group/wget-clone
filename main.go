package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"wget/pkg"
)

var (
	B         bool
	O         string
	P         string
	I         string
	RateLimit string
	Mirror    bool
	Reject    string
	Exclude   string
)

func init() {
	flag.BoolVar(&B, "B", false, "Output in wget-log")
	flag.StringVar(&O, "O", "", "Specify download filename")
	flag.StringVar(&P, "P", "./", "Specify download directory")
	flag.StringVar(&I, "i", "", "Download multiple files")
	flag.StringVar(&RateLimit, "rate-limit", "10000M", "The rate limit in k = KB/s or  M = MB/s")
	flag.BoolVar(&Mirror, "mirror", false, "Mirror the whole site")
	flag.StringVar(&Reject, "reject", "", "Reject files")
	flag.StringVar(&Reject, "R", "", "Reject files")
	flag.StringVar(&Exclude, "exclude", "", "Exclude directory")
	flag.StringVar(&Exclude, "X", "", "Exclude directory")
}

func main() {
	flag.Parse()
	url := flag.Arg(0)
	
	rate, err := pkg.GetRateLimit(RateLimit)
	if err != nil {
		fmt.Println(err)
		return
	}

	if Mirror {
		pkg.Mirror(url, Exclude, Reject)
		
	} else if I != "" && (O == "") {
		urls, err := pkg.ReadDownloadFile(I)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, fileNames, err := pkg.DownloadMultipleFiles(P, urls, rate)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		filePath := P
		if strings.Contains(P, "~") {
			usr, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				return
			}
			filePath = path.Join(usr, P[1:])
		}

		for i, r := range resp {
			pkg.SaveBytesToFile(filePath, fileNames[i], r)
			fmt.Printf("Downloaded file %s \n", fileNames[i])
		}
	} else {
		fileName := path.Base(url) // extract the file name from the url
		if O != "" && I == "" {
			fileName = O
		}
		// Send GET request to the provided URL
		response, err := http.Get(url)
		if err != nil {
			// Return nil and error if request fails
			return
		}
		defer response.Body.Close()
		download := &pkg.Download{Response: response, StartTime: time.Now(), ContentLength: float64(response.ContentLength), BarWidth: pkg.GetTerminalLength(), Url: url}
		if P != "./" {
			download.Path = P + "/" + fileName
		}
		download.PrintOrLogBefore(B)
		resp, err := download.DownloadFile(response, rate)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		filePath := P
		if strings.Contains(P, "~") {
			usr, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				return
			}
			filePath = path.Join(usr, P[1:])
		}
		pkg.SaveBytesToFile(filePath, fileName, resp)
		download.PrintOrLogAfter(B)

	}
}
