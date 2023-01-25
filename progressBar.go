package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Download struct {
	response          *http.Response
	contentLength     float64
	currentBytes      float64
	startTime         time.Time
	percentage        int
	previousBarLength int
}

func (data *Download) StartProgressBar() {
	// Create Empty Progress Bar
	fmt.Println(data.CreateBarString())
}

func (data *Download) CreateBarString() string {
	return ByteToUnit(data.currentBytes) + " / " + ByteToUnit(data.contentLength) + " [>                                                                      ] " + strconv.Itoa(data.percentage) + "% " + data.RateOfDownload() + " " + data.TimeRemaining()
}

// This function takes in a int representing bytes and returns a string of the input in the appropriate unit
func ByteToUnit(byteCount float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unit := 0
	for byteCount >= 1024 && unit < 4 {
		byteCount /= 1024
		unit++
	}
	return strconv.Itoa(int(byteCount)) + " " + units[unit]
}

func (data *Download) RateOfDownload() string {
	elapsed := time.Now().Sub(data.startTime)
	return ByteToUnit(data.currentBytes/elapsed.Seconds()) + "/s"
}

func (data *Download) TimeRemaining() string {
	BytesPerSecond := data.currentBytes / time.Now().Sub(data.startTime).Seconds()
	RemainingBytes := (data.contentLength - data.currentBytes)
	RemainingTime := RemainingBytes / BytesPerSecond
	return strconv.Itoa(int(RemainingTime)) + "s"
}
