package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
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

func DownloadMultipleFiles(P string, urls []string, ratelimit int64) {
	var wg sync.WaitGroup
	for _, URL := range urls {
		go func(URL string) {
			wg.Add(1)
			defer wg.Done()
			fileName := path.Base(URL)

			client := &http.Client{}

			res, err := http.NewRequest("GET", URL, nil)
			if err != nil {
				return
			}

			res.Header.Set("User-Agent", "golang-wget-project")
			response, err := client.Do(res)
			if err != nil {
				return
			}
			fmt.Println("Download started --> ", fileName)
			defer response.Body.Close()
			download := &Download{Response: response, StartTime: time.Now(), ContentLength: float64(response.ContentLength), BarWidth: GetTerminalLength(), Url: URL}
			if P != "./" {
				download.Path = P + "/" + fileName
			}

			download.HideBar = true

			b, err := download.DownloadFile(response, ratelimit)
			if err != nil {
				fmt.Println("Error downloading ", fileName, " ", err)
				return
			}
			SaveBytesToFile(P, fileName, b)
			fmt.Printf("Downloaded file %s \n", fileName)
		}(URL)
	}
	wg.Wait()
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
