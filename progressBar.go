package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Download struct {
	response          *http.Response
	contentLength     float64
	currentBytes      float64
	startTime         time.Time
	percentage        int
	previousBarLength int
	ProgressBar       string
	BarWidth          int
}

func (data *Download) UpdateProgressBar() {
	buffer := make([]byte, 100000000)
	for {
		r, err := data.response.Body.Read(buffer)
		if err != nil {
			return
		}
		data.currentBytes += float64(r)
		data.percentage = int((data.currentBytes * 100) / data.contentLength)
		fmt.Fprintf(os.Stdout, "\r")
		data.PrintProgressBar()
	}
}

func (data *Download) StartProgressBar() {
	// Create Empty Progress Bar
	data.UpdateProgressBar()
}

func (data *Download) PrintProgressBar() {
	data.CreateProgressBar()
	fmt.Print(data.ProgressBar)
}

func (data *Download) CreateProgressBar() {
	data.ProgressBar = ByteToUnit(data.currentBytes) + " / " + ByteToUnit(data.contentLength) + data.ProgressString() + strconv.Itoa(data.percentage) + "% " + data.RateOfDownload() + " " + data.TimeRemaining()
	data.CheckNewLengthWithPrevious()
}

// This function takes in a int representing bytes and returns a string of the input in the appropriate unit
func ByteToUnit(byteCount float64) string {
	units := []string{"B", "MB", "MB", "GB", "TB"}
	unit := 0
	for byteCount > 1024 && unit < 4 {
		byteCount /= 1024
		unit++
	}
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

func (data *Download) ProgressString() string {
	var s string = " ["
	n := float64(data.BarWidth) * (float64(data.percentage) / 100)

	for i := 0; i < data.BarWidth; i++ {
		if i == int(n) {
			s += ">"
		} else if i < int(n) {
			s += "="
		} else {
			s += " "
		}
	}
	s += "] "
	return s
}

func (data *Download) CheckNewLengthWithPrevious() {
	if len(data.ProgressBar) < data.previousBarLength {
		spacebuffer := strings.Repeat(" ", data.previousBarLength-len(data.ProgressBar))
		data.ProgressBar += spacebuffer
	}
	data.previousBarLength = len(data.ProgressBar)
}

func GetTerminalLength() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Println("Error getting terminal size: ", err)
		return 50
	}
	T := strings.Fields(strings.TrimSpace(string(out)))
	wid := T[1]
	w, convErr := strconv.Atoi(wid)
	if convErr != nil {
		log.Println(err)
		return 50
	}
	return w / 5
}
