package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadFile(URL string, rateLimit int64) ([]byte, error) {
	// Send GET request to the provided URL
	response, err := http.Get(URL)
	if err != nil {
		// Return nil and error if request fails
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		// Return nil and error if response status is not OK
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	if rateLimit > 0 {
		// Create a ticker with a 1 second interval
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		buf := make([]byte, rateLimit)
		for {
			select {
			case <-ticker.C:
				// Read up to rateLimit bytes from the response body
				n, err := response.Body.Read(buf)
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
	} else {
		// Read from the response body without rate limiting
		_, err = io.Copy(&data, response.Body)
		if err != nil {
			// Return nil and error if an error occurred
			return nil, err
		}
	}
	// Return the bytes and a nil error if successful
	return data.Bytes(), nil
}

func DownloadMultipleFiles(urls []string, ratelimit int64) ([][]byte, error) {
	done := make(chan []byte, len(urls))
	errch := make(chan error, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			b, err := DownloadFile(URL, ratelimit)
			if err != nil {
				errch <- err
				done <- nil
				return
			}
			done <- b
			errch <- nil
		}(URL)
	}
	bytesArray := make([][]byte, 0)
	var errStr string
	for i := 0; i < len(urls); i++ {
		bytesArray = append(bytesArray, <-done)
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}
	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}
	return bytesArray, err
}

func SaveBytesToFile(fileName string, r []byte) {
	err := os.WriteFile(fileName, r, 0o644)
	if err != nil {
		fmt.Println(err)
	}
}
