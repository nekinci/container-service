package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
)

func NewNoAvailable() http.Response {
	f, _ := os.ReadFile("./resources/no-available.html")

	response := http.Response{
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewReader(f)),
		ContentLength: int64(len(f)),
	}

	return response
}

func NewNoLongerAvailable() http.Response{
	f, _ := os.ReadFile("./resources/no-longer-available.html")

	response := http.Response{
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewReader(f)),
		ContentLength: int64(len(f)),
	}

	return response
}

func IsEmpty(byteArr []byte) bool {

	for b := range byteArr {
		if b != 0 {
			return false
		}
	}

	return true
}