package filemail

import (
	"encoding/json"
	"errors"
	"main/utils"
	"strconv"
)

const (
	apiBase = "https://www.filemail.com/api/transfer/"
	referer = "https://www.filemail.com/"
)

func initUpload(size int64, headers map[string]string) (*Initialize, error) {
	initUrl := apiBase + "initialize"
	params := map[string]string{
		"sourcedetails": "fupload4.0 @ https://www.filemail.com/",
		"days":          "7",
		"confirmation":  "true",
		"transfersize":  strconv.FormatInt(size, 10),
	}
	respBody, err := utils.DoPost(initUrl, params, headers, nil)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()
	var obj Initialize
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return nil, err
	}
	if obj.Responsestatus != "OK" {
		return nil, errors.New("Bad response.")
	}
	return &obj, nil
}

func upload(uploadUrl, path string, size, byteLimit int64, params, headers map[string]string) error {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file", size, byteLimit, nil, params, headers)
	if err != nil {
		return err
	}
	respBody.Close()
	return nil
}

func finalizeUpload(params, headers map[string]string) (string, error) {
	finalizeUrl := apiBase + "complete"
	params["failed"] = "false"
	respBody, err := utils.DoPost(finalizeUrl, params, headers, nil)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj Finalize
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	if obj.Responsestatus != "OK" {
		return "", errors.New("Bad response.")
	}
	return obj.Downloadurl, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "5GB")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": referer,
		"Source":  "Web",
	}
	initMeta, err := initUpload(size, headers)
	if err != nil {
		return "", err
	}
	params := map[string]string{
		"transferid":  initMeta.Transferid,
		"transferkey": initMeta.Transferkey,
	}
	err = upload(initMeta.Transferurl, path, size, args.ByteLimit, params, headers)
	if err != nil {
		return "", err
	}
	fileUrl, err := finalizeUpload(params, headers)
	return fileUrl, err
}
