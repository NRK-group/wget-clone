package pkg

import (
	"fmt"
	"strings"
	"time"
)

func (data *Download) PrintBefore() {
	fmt.Println("start at " + strings.Split(data.StartTime.String(), ".")[0])
	fmt.Println("sending request, awaiting response... " + data.Response.Status)
	fmt.Println("content size: " + fmt.Sprintf("%0.0f", data.ContentLength) + " [~" + MegOrGig(data.ContentLength) + "]")
	fmt.Println("saving file to: " + data.Path)
}

func (data *Download) PrintAfter() {
	fmt.Println("\n\nDownloaded [" + data.Url + "]")
	fmt.Println("finished at ", strings.Split(time.Now().String(), ".")[0])
}
