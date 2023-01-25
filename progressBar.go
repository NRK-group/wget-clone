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

func (data *Download) UpdateProgressBar() {
	buffer := make([]byte, 100000000)
	for {
		r, err := data.response.Body.Read(buffer)
		if err != nil {
			return
		}
		data.currentBytes += float64(r)
		data.percentage = int((data.currentBytes*100) / data.contentLength)
		// fmt.Println("New PERCENTAGE === ", ((int(data.contentLength) - int(data.currentBytes)) / int(data.contentLength)) * 100)
		fmt.Println(data.CreateBarString())
	}
}

func (data *Download) StartProgressBar() {
	// Create Empty Progress Bar
	fmt.Println(data.CreateBarString())
	data.UpdateProgressBar()
}

func (data *Download) CreateBarString() string {
	return ByteToUnit(data.currentBytes) + " / " + ByteToUnit(data.contentLength) + " [] " + strconv.Itoa(data.percentage) + "% " + data.RateOfDownload() + " " + data.TimeRemaining()
}

// This function takes in a int representing bytes and returns a string of the input in the appropriate unit
func ByteToUnit(byteCount float64) string {
	units := []string{"B", "MB", "MB", "GB", "TB"}
	unit := 0
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()

	// fmt.Println("BYTES BEFORE: ", byteCount)
	for byteCount > 1024 && unit < 4 {
		byteCount /= 1024
		unit++
	}
	// fmt.Println("BYTES After: ", byteCount)
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()


	return strconv.Itoa(int(byteCount)) + units[unit]
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
