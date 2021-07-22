package nasa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type NasaContext struct {
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
func (ctx *NasaContext) GetData(w http.ResponseWriter, r *http.Request) {

	fileName := time.Now().Format("2006-01-02")
	b, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		fileData, ftype, httpErr := ctx.getFileFromHttp()
		if httpErr == nil {
			w.Header().Set("Content-Type", ftype+"; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fileData)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.Header().Set("Content-Type", "image/png; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}

}

func (ctx *NasaContext) readDataInfo(data []byte) (*NasaData, error) {
	fmt.Printf("body is %s\n", string(data))
	var info NasaData
	err := json.Unmarshal(data, &info)
	if err != nil {
		fmt.Printf("err is %s\n", err.Error())
		return nil, errors.New("Failure")
	}
	fmt.Printf("file %s located at %s\n", info.Date, info.Url)
	return &info, nil
}

func (ctx *NasaContext) getImage(url string) ([]byte, error) {
	picture, err := http.Get(url)
	if err != nil {
		fmt.Printf("Cannot get picture url\n")
		return nil, errors.New("Failure")
	}

	file, err := ioutil.ReadAll(picture.Body)
	if err != nil {
		fmt.Printf("Cannot read picture body")
		return nil, errors.New("Failure")
	}

	return file, nil
}

func (ctx *NasaContext) getFileFromHttp() ([]byte, string, error) {
	resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=" + ctx.ApiKey)
	if err != nil {
		fmt.Printf("error occured cannot get file, http return error %s", err.Error())
		return nil, "", errors.New("Failure")
	}
	fmt.Printf("Http call status is %s\n", resp.Status)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Cannot read body")
		return nil, "", errors.New("Failure")
	}

	info, err := ctx.readDataInfo(data)
	if err != nil {
		fmt.Printf("Cannot unmarshall body")
		return nil, "", errors.New("Failure")
	}

	picture, err := ctx.getImage(info.Url)
	if err != nil {
		fmt.Printf("Cannot get picture url\n")
		return nil, "", errors.New("Failure")
	}
	os.WriteFile("./"+info.Date, picture, 0666)

	fileType := "application/octet-stream"

	switch strings.Split(info.Url, ".")[len(strings.Split(info.Url, "."))-1] {
	case "png":
		fileType = "image/png"
	case "jpg":
		fileType = "image/jpg"
	}
	return picture, fileType, nil
}
