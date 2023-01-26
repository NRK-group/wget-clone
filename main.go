package main

import (
	"net/http"
	"time"
)

func main() {
	response, o := http.Get("https://golang.org/dl/go1.16.3.linux-amd64.tar.gz")
	if o != nil {
		return
	}
	currentDownload := Download{startTime: time.Now(), contentLength: float64(response.ContentLength), response: response, BarWidth: GetTerminalLength()}
	// fmt.Println()
	// fmt.Print("D")

	// fmt.Print("Hello,	")
	// fmt.Print("world!")
	// fmt.Println()
	// fmt.Fprintf(os.Stdout, "\a")

	// fmt.Print("Goodbye, ")

	currentDownload.StartProgressBar()
}
