package catbox

import (
	"io"
	"main/utils"
)

func upload(uploadUrl, path string, size int64, formMap, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "fileToUpload", size, formMap, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	bodyBytes, err := io.ReadAll(respBody)
	return string(bodyBytes), err
}

func Run(args *utils.Args, path string) (string, error) {
	uploadUrl := "https://catbox.moe/user/api.php"
	size, err := utils.CheckSize(path, "200MB")
	if err != nil {
		return "", err
	}
	formMap := map[string]string{
		"userhash": "",
		"reqtype":  "fileupload",
	}
	headers := map[string]string{
		"Referer": "https://catbox.moe/",
	}
	fileUrl, err := upload(uploadUrl, path, size, formMap, headers)
	return fileUrl, err
}
