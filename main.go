package main

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

type ProgressBar struct {
	startTime  time.Time
	endTime    time.Time
	width      int
	contentLen int64
	current    int
	percentage int
}

func (bar *ProgressBar) UpdateTerminalBar(current int) {
	s := ""
	percentage := math.Round((float64(current) / float64(bar.contentLen)) * 100)
	s += ("[")
	n := int((percentage / 100) * float64(bar.width))

	// for i := 0; i < bar.width; i++ {
	// 	fmt.Print("\033[F")
	// }
	fmt.Printf("\033[A\033[2K")
	// fmt.Printf("\033[")

	for i := 0; i < bar.width; i++ {
		if i <= n {
			s += (">")
		} else {
			s += (" ")
		}
	}
	s += ("]")
	s += (fmt.Sprintf(" %v%%", percentage))
	fmt.Println("Progress: ", s)
}

func main() {
	buf := make([]byte, 1000000)
	var bar ProgressBar

	response, o := http.Get("https://www.bigdataframework.org/wp-content/uploads/2019/11/2.jpg")
	if o != nil {
		fmt.Println("Error :", o)
		return
	}
	fmt.Print("Request content length: ")
	bar.contentLen = response.ContentLength
	bar.width = 70
	fmt.Println(response.ContentLength)
	fmt.Println()

	// var writer io.Writer
	var b int
	go func() {
		for {

			y := 0
			for i := 0; i < len(buf); i++ {
				if buf[i] != 0 {
					y++
				}
			}
			// fmt.Println()
			// fmt.Println()
			// fmt.Println()
			// fmt.Println("Index of buff that isn't 0 === ", y)
			// fmt.Println()
			// fmt.Println()
			// fmt.Println()
			// Read stores the current position in an internal pointer and carries on from the next iteration
			r, err := response.Body.Read(buf)
			if err != nil {
				fmt.Println("\n Error: ", err)
				return
			}
			// if b != r {
			// fmt.Println()
			// fmt.Println()
			// fmt.Println()
			// fmt.Println("before === ", b)
			// fmt.Println()
			// fmt.Println()
			// fmt.Println("after ==== ", r+b)
			// fmt.Println()
			// fmt.Println()
			// fmt.Println()
			// fmt.Println("Print b BEFORE == ", b)
			// fmt.Println()
			// fmt.Println()
			b += r
			// fmt.Println("Print b AFTER == ", b)
			bar.UpdateTerminalBar(b)
			// }

		}
	}()
	time.Sleep(5 * time.Second)
	fmt.Println("FINISHED")
}
