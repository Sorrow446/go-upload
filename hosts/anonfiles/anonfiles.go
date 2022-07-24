package anonfiles

import (
	"encoding/json"
	"errors"
	"main/utils"
)

const (
	uploadUrl = "https://api.anonfiles.com/upload"
	referer   = "https://anonfiles.com/"
)

func upload(path string, size, byteLimit int64, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file", size, byteLimit, nil, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
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
	size, err := utils.CheckSize(path, "20GB")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": referer,
	}
	fileUrl, err := upload(path, size, args.ByteLimit, headers)
	return fileUrl, err
}
