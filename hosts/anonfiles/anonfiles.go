package anonfiles

import (
	"encoding/json"
	"errors"
	"main/utils"
)

func upload(uploadUrl, path string, size int64, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file", size, nil, nil, nil)
	if err != nil {
		return "", err
	}
	if respBody != nil {
		defer respBody.Close()
	}
	var obj Upload
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	if !obj.Status {
		return "", errors.New("Bad response. " + obj.Error.Type)
	} else if obj.Data.File.Metadata.Size.Bytes != size {
		return "", errors.New("Byte count mismatch.")
	}
	return obj.Data.File.URL.Full, nil
}

func Run(args *utils.Args, path string) (string, error) {
	uploadUrl := "https://api.anonfiles.com/upload"
	size, err := utils.CheckSize(path, "20GB")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": "https://anonfiles.com/",
	}
	fileUrl, err := upload(uploadUrl, path, size, headers)
	return fileUrl, err
}
