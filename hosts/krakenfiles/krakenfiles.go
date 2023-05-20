package krakenfiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/utils"
)

const (
	referer          = "https://krakenfiles.com/"
	uloadUrlRegexStr = `url: "([^"]+)"`
)

func getUploadUrl() (string, error) {
	match, err := utils.FindHtmlSubmatch(referer, uloadUrlRegexStr)
	if err != nil {
		return "", err
	}
	if match == nil {
		return "", errors.New("No regex match.")
	}
	return match[1], nil
}
func upload(uploadUrl, path string, size, ByteLimit int64, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "files[]", size, ByteLimit, nil, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj Upload
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	file := obj.Files[0]
	if file.Error != "" {
		return "", errors.New("Bad response: " + file.Error)
	}
	return referer + file.URL[1:], nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "1GB")
	if err != nil {
		return "", err
	}
	uploadUrl, err := getUploadUrl()
	if err != nil {
		fmt.Println("Failed to get upload URL.")
		return "", err
	}
	headers := map[string]string{
		"Referer": referer,
	}
	fileUrl, err := upload(uploadUrl, path, size, args.ByteLimit, headers)
	return fileUrl, err
}
