package gofile

import (
	"encoding/json"
	"errors"
	"main/utils"
)

func getServer() (string, error) {
	serverUrl := "https://api.gofile.io/getServer"
	respBody, err := utils.DoGet(serverUrl, nil, nil)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj GetServer
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	if obj.Status != "ok" {
		return "", errors.New("Bad response.")
	}
	return obj.Data.Server, nil
}

func upload(uploadUrl, path string, size int64, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file", size, nil, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj Upload
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	if obj.Status != "ok" {
		return "", errors.New("Bad response.")
	}
	return obj.Data.DownloadPage, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "unlim")
	if err != nil {
		return "", err
	}
	server, err := getServer()
	if err != nil {
		return "", err
	}
	uploadUrl := "https://" + server + ".gofile.io/uploadFile"
	headers := map[string]string{
		"Referer": "https://gofile.io/",
	}
	url, err := upload(uploadUrl, path, size, headers)
	return url, err
}
