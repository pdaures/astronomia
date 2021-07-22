package nasa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	ApiKey string
}

type NasaData struct {
	Copyright       string
	Date            string
	Explanation     string
	Hdurl           string
	Media_type      string
	Service_version string
	Title           string
	Url             string
}

// function that will serve file if exist or download it if not
func Handler(apiKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// file not in local cache, let's get it from NASA
		b, mimeType, err := fetchTodayImage(apiKey)
		if err == nil {
			w.Header().Set("Content-Type", mimeType)
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(b); err != nil {
				fmt.Printf("failed to send response to client, %+v\n", err)
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func fetchTodayImage(apiKey string) ([]byte, string, error) {
	url := "https://api.nasa.gov/planetary/apod?api_key=" + apiKey
	data, err := get(url)
	if err != nil {
		return nil, "", err
	}

	var info NasaData
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, "", fmt.Errorf("failed to parse nasa data response, %+v", err)
	}

	picture, err := get(info.Url)
	if err != nil {
		return nil, "", err
	}

	fileType := "application/octet-stream"

	switch strings.Split(info.Url, ".")[len(strings.Split(info.Url, "."))-1] {
	case "png":
		fileType = "image/png"
	case "jpg":
		fileType = "image/jpg"
	}
	return picture, fileType, nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %s failed, %+v\n", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s returned %s", url, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from GET %s, %+v", url, err)
	}

	return data, nil
}
