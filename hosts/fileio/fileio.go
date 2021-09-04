package fileio

import (
	"encoding/json"
	"errors"
	"main/utils"
)

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
	if !obj.Success {
		return "", errors.New("Bad response.")
	} else if obj.Size != size {
		return "", errors.New("Byte count mismatch.")
	}
	return obj.Link, nil
}

func Run(args *utils.Args, path string) (string, error) {
	uploadUrl := "https://file.io/"
	size, err := utils.CheckSize(path, "100MB")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": "https://www.file.io/",
	}
	fileUrl, err := upload(uploadUrl, path, size, headers)
	return fileUrl, err
}
