package pkg

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Download struct {
	Response          *http.Response
	ContentLength     float64
	CurrentBytes      float64
	StartTime         time.Time
	Percentage        int
	PreviousBarLength int
	ProgressBar       string
	BarWidth          int
}

func (data *Download) UpdateProgressBar(n int) {
	
		data.CurrentBytes += float64(n)
		data.Percentage = int((data.CurrentBytes * 100) / data.ContentLength)
		fmt.Fprintf(os.Stdout, "\r")
		data.PrintProgressBar()
}

func (data *Download) PrintProgressBar() {
	data.CreateProgressBar()
	fmt.Print(data.ProgressBar)
}

func (data *Download) CreateProgressBar() {
	data.ProgressBar = ByteToUnit(data.CurrentBytes) + " / " + ByteToUnit(data.ContentLength) + data.ProgressString() + strconv.Itoa(data.Percentage) + "% " + data.RateOfCurrent() + " " + data.TimeRemaining()
	data.CheckNewLengthWithPrevious()
}

// This function takes in a int representing bytes and returns a string of the input in the appropriate unit
func ByteToUnit(byteCount float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unit := 0
	for byteCount > 1024 && unit < 4 {
		byteCount /= 1024
		unit++
	}
	return strconv.Itoa(int(byteCount)) + units[unit]
}

func (data *Download) RateOfCurrent() string {
	elapsed := time.Now().Sub(data.StartTime)
	return ByteToUnit(data.CurrentBytes/elapsed.Seconds()) + "/s"
}

func (data *Download) TimeRemaining() string {
	BytesPerSecond := data.CurrentBytes / time.Now().Sub(data.StartTime).Seconds()
	RemainingBytes := (data.ContentLength - data.CurrentBytes)
	RemainingTime := RemainingBytes / BytesPerSecond
	return strconv.Itoa(int(RemainingTime)) + "s"
}

func (data *Download) ProgressString() string {
	var s string = " ["
	n := float64(data.BarWidth) * (float64(data.Percentage) / 100)

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
	if len(data.ProgressBar) < data.PreviousBarLength {
		spacebuffer := strings.Repeat(" ", data.PreviousBarLength-len(data.ProgressBar))
		data.ProgressBar += spacebuffer
	}
	data.PreviousBarLength = len(data.ProgressBar)
}

func GetTerminalLength() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 50
	}
	T := strings.Fields(strings.TrimSpace(string(out)))
	wid := T[1]
	w, convErr := strconv.Atoi(wid)
	if convErr != nil {
		return 50
	}
	return w / 5
}
