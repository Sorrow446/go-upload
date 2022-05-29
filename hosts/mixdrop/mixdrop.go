package mixdrop

import (
	"encoding/json"
	"main/utils"
)

const (
	referer   = "https://mixdrop.co/"
	uploadUrl = "https://ul.mixdrop.co/up"
)

func upload(path string, size, byteLimit int64, headers, formMap map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "files", size, byteLimit, formMap, nil, nil)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj UploadResp
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	return referer + "f/" + obj.File.Ref, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "unlim")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": referer,
	}
	formMap := map[string]string{
		"upload": "1",
	}
	fileUrl, err := upload(path, size, args.ByteLimit, headers, formMap)
	return fileUrl, err
}
