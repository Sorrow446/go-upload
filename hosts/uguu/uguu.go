package uguu

import (
	"encoding/json"
	"errors"
	"main/utils"
	"strconv"
)

const (
	referer   = "https://uguu.se/"
	uploadUrl = referer + "upload.php"
)

func upload(path string, size, ByteLimit int64, headers map[string]string) (string, error) {
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
	retSize, err := strconv.ParseInt(obj.Files[0].Size, 10, 64)
	if err != nil {
		return "", err
	}
	if !obj.Success {
		return "", errors.New("Bad response.")
	} else if retSize != size {
		return "", errors.New("Byte count mismatch.")
	}
	return obj.Files[0].URL, nil

}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "128MB")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": referer,
	}
	fileUrl, err := upload(path, size, args.ByteLimit, headers)
	return fileUrl, err
}
