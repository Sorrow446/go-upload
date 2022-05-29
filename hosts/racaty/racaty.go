package racaty

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/utils"
)

const (
	referer          = "https://racaty.net/"
	uloadUrlRegexStr = `<form id="uploadfile" action="([^"]+)"`
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

func upload(uploadUrl, path string, size, byteLimit int64, headers, formMap map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file_0", size, byteLimit, formMap, nil, nil)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj UploadResp
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	file := obj[0]
	if file.FileStatus != "OK" {
		return "", errors.New("Bad response.")
	}
	return referer + file.FileCode, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "unlim")
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
	formMap := map[string]string{
		"sess_id":   "",
		"utype":     "anon",
		"link_rcpt": "",
		"link_pass": "",
		"to_folder": "",
		"keepalive": "1",
	}
	fileUrl, err := upload(uploadUrl, path, size, args.ByteLimit, headers, formMap)
	return fileUrl, err
}
