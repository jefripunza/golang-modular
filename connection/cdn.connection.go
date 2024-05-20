package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"project/env"
	"strings"
)

type CDN struct{}

func (ref CDN) Upload(filePath string, bucketName string) (string, error) {
	cdn_host_url := env.GetCdnHostUrl()
	url := cdn_host_url + "/upload/multipart"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	_ = writer.WriteField("bucket_name", bucketName)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("x-api-key", "keynya")
	req.Header.Add("x-api-pass", "passnya")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// Parsing respons JSON
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Mengambil nilai dari kunci "file_name"
	fileName, ok := response["file_name"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	fileName = cdn_host_url + "/file/" + bucketName + "/" + fileName
	fileName = strings.ReplaceAll(fileName, "//", "/")
	return fileName, nil
}
