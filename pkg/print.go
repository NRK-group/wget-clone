package pkg

import (
	"fmt"
	"os"
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

func (data *Download) WriteToLogBefore() {
	data.HideBar = true
	fmt.Println("Output will be written to ‘wget-log’.")
	var s string
	// start at 2017-10-14 03:46:06
	s += "start at " + strings.Split(data.StartTime.String(), ".")[0] + "\n"
	// sending request, awaiting response... status 200 OK
	s += "sending request, awaiting response..." + data.Response.Status + "\n"
	// content size: 56370 [~0.06MB]
	s += "content size: " + fmt.Sprintf("%0.0f", data.ContentLength) + " [~" + MegOrGig(data.ContentLength) + "]" + "\n"
	// saving file to: ./EMtmPFLWkAA8CIS.jpg
	s += "saving file to: " + data.Path + "\n"

	log, err := os.OpenFile("wget-log", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error writing to log file: ", err)
	}
	log.Write([]byte(s))
}

func (data *Download) WriteToLogAfter() {
	var s string
	s += "Downloaded [" + data.Url + "]" + "\n"
	s += "finished at " + strings.Split(time.Now().String(), ".")[0] + "\n"
	log, err := os.OpenFile("wget-log", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error writing to log file: ", err)
	}
	log.Write([]byte(s))
}
