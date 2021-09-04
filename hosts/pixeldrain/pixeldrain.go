package pixeldrain

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
	}
	url := "https://pixeldrain.com/u/" + obj.ID
	return url, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "5GB")
	if err != nil {
		return "", err
	}
	uploadUrl := "https://pixeldrain.com/api/file"
	headers := map[string]string{
		"Referer": "https://pixeldrain.com/",
	}
	fileUrl, err := upload(uploadUrl, path, size, headers)
	return fileUrl, err
}
