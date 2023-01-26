package main

import (
	"flag"
	"fmt"

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
	flag.StringVar(&RateLimit, "rate-limit", "0", "The rate limit in k = KB/s or  M = MB/s")
	flag.BoolVar(&Mirror, "mirror", false, "Mirror the whole site")
	flag.StringVar(&Reject, "reject", "", "Reject files")
	flag.StringVar(&Reject, "R", "", "Reject files")
	flag.StringVar(&Exclude, "exclude", "", "Exclude directory")
	flag.StringVar(&Exclude, "X", "", "Exclude directory")
}

func main() {
	flag.Parse()
	if B {
		fmt.Println("Output in wget-log is enabled")
	} else {
		fmt.Println("Output in wget-log is disabled")
	}
	fmt.Println("filename:", O)
	fmt.Println("directory:", P)
	if I {
		fmt.Println("Download multiple files is enabled")
	} else {
		fmt.Println("Download multiple files is disabled")
	}
	fmt.Println("Rate Limit: ", RateLimit)
	if Mirror {
		fmt.Println("Mirror the whole site is enabled")
		url := flag.Arg(0)
		pkg.Mirror(url, Exclude, Reject)
	} else {
		fmt.Println("Mirror the whole site is disabled")
	}
	/*
		fmt.Println("Reject:", Reject)
		fmt.Println("Exclude:", Exclude)
		fmt.Println("URL:", flag.Arg(0))
		url := flag.Arg(0)
		rate, err := pkg.GetRateLimit(RateLimit)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := pkg.DownloadFile(url, rate)
		if err != nil {
			fmt.Println(err)
			return
		}
		fileName := path.Base(url) // extract the file name from the url
		pkg.SaveBytesToFile("./js/"+fileName, resp)
	*/
}
