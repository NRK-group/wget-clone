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
	I         bool
	RateLimit string
	Mirror    bool
	Reject    string
	Exclude   string
)

func init() {
	flag.BoolVar(&B, "B", false, "Output in wget-log")
	flag.StringVar(&O, "O", "", "Specify download filename")
	flag.StringVar(&P, "P", "./", "Specify download directory")
	flag.BoolVar(&I, "i", false, "Download multiple files")
	flag.StringVar(&RateLimit, "rate-limit", "10000M", "The rate limit in k = KB/s or  M = MB/s")
	flag.BoolVar(&Mirror, "mirror", false, "Mirror the whole site")
	flag.StringVar(&Reject, "reject, R", "", "Reject files")
	flag.StringVar(&Exclude, "exclude, X", "", "Exclude directory")
}

func main() {
	flag.Parse()
	url := flag.Arg(0)
	fileName := path.Base(url) // extract the file name from the url
	if O != "" && !I {
		fileName = O
	}
	rate, err := pkg.GetRateLimit(RateLimit)
	if err != nil {
		fmt.Println(err)
		return
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
	if I && (O != "") {
		download.HideBar = true
		fmt.Println("Download multiple files is enabled")
	} else {

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
