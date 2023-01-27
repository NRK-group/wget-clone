package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func (download *Download) DownloadFile(response *http.Response, rateLimit int64) ([]byte, error) {
	if response.StatusCode != http.StatusOK {
		// Return nil and error if response status is not OK
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	// Create a ticker with a 1 second interval
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	buf := make([]byte, rateLimit)
	for {
		select {
		case <-ticker.C:
			// Read up to rateLimit bytes from the response body
			n, err := response.Body.Read(buf)
			download.UpdateProgressBar(n)
			if err != nil {
				if err == io.EOF {
					// Return the bytes and a nil error if EOF is reached
					return data.Bytes(), nil
				}
				// Return nil and error if an error occurred
				return nil, err
			}
			if _, err := data.Write(buf[:n]); err != nil {
				// Return nil and error if an error occurred
				return nil, err
			}
		}
	}
}

func DownloadMultipleFiles(P string, urls []string, ratelimit int64) ([][]byte, []string, error) {
	done := make(chan []byte, len(urls))
	errch := make(chan error, len(urls))
	fileNames := make(chan string, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			fileName := path.Base(URL)

			// Send GET request to the provided URL
			response, err := http.Get(URL)
			if err != nil {
				// Return nil and error if request fails
				return
			}
			defer response.Body.Close()
			download := &Download{Response: response, StartTime: time.Now(), ContentLength: float64(response.ContentLength), BarWidth: GetTerminalLength(), Url: URL, }
			if P != "./" {
				download.Path = P + "/" + fileName
			}

			download.HideBar = true
			fmt.Println("Download Started: ", URL)

			b, err := download.DownloadFile(response, ratelimit)
			if err != nil {
				errch <- err
				done <- nil
				fileNames <- ""
				return
			}
			SaveBytesToFile(P, fileName, b)
			// done <- b
			errch <- nil
			fileNames <- fileName
		}(URL)
	}
	bytesArray := make([][]byte, 0)
	namesArray := make([]string, 0)
	var errStr string
	for i := 0; i < len(urls); i++ {
		bytesArray = append(bytesArray, <-done)
		namesArray = append(namesArray, <-fileNames)
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}
	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}
	return bytesArray, namesArray, err
}

func SaveBytesToFile(pathName, fileName string, r []byte) {
	if pathName != "./" {
		err := os.MkdirAll(pathName, 0o755)
		if err != nil {
			// handle the error
			fmt.Println(err)
			return
		}
		pathName += "/"
	}
	err := os.WriteFile(pathName+fileName, r, 0o644)
	if err != nil {
		fmt.Println(err)
	}
}
