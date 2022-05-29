package transfersh

import (
	"io/ioutil"
	"main/utils"
	"strings"
)

const uploadUrl = "https://transfer.sh/"

func upload(uploadUrl, path string, size, byteLimit int64) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file", size, byteLimit, nil, nil, nil)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	bodyBytes, err := ioutil.ReadAll(respBody)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bodyBytes)), nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "unlim")
	if err != nil {
		return "", err
	}
	fileUrl, err := upload(uploadUrl, path, size, args.ByteLimit)
	return fileUrl, err
}
