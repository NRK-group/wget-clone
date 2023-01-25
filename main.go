package main

import (
	"time"
)

func main() {
	// response, o := http.Get("https://golang.org/dl/go1.16.3.linux-amd64.tar.gz")
	// if o != nil {
	// return
	// }
	// ShowProgressBar(response)
	// s := []string{"1","22", "333", "4444", "55555"}
	// for i:= 0; i< 5 ;i++ {
	// 		fmt.Print(s[i])
	// 		// bar.Increment()
	// 		fmt.Fprintf(os.Stdout, "\x1B[K\r")
	// 		time.Sleep(time.Second*2)
	// }

	currentDownload := Download{startTime: time.Now(), currentBytes: 1024, contentLength: 1024000}
	currentDownload.StartProgressBar()
}
