package letsupload

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/utils"
	"strconv"
	"time"
)

const (
	referer           = "https://letsupload.io/"
	uploaderUrl       = referer + "assets/js/uploader.js"
	uploadUrlRegexStr = `url: '([^']+)'`
	formDataRegexStr  = `data.formData = {_sessionid: '([^']+)', cTracker: '([^']+)'`
)

func getEpochStr() string {
	epoch := strconv.FormatInt(time.Now().Unix(), 10)
	return epoch
}

func extractMeta() (*Meta, error) {
	epoch := getEpochStr()
	html, err := utils.GetHtml(uploaderUrl + "?r=" + epoch)
	if err != nil {
		return nil, err
	}
	uploadUrlMatch := utils.FindStringSubmatch(html, uploadUrlRegexStr)
	if uploadUrlMatch == nil {
		return nil, errors.New("No regex match for upload URL.")
	}
	formDataMatch := utils.FindStringSubmatch(html, formDataRegexStr)
	if formDataMatch == nil {
		return nil, errors.New("No regex match for form data.")
	}
	meta := &Meta{
		UploadURL: uploadUrlMatch[1],
		SessionID: formDataMatch[1],
		Tracker:   formDataMatch[2],
	}
	return meta, nil
}

func upload(uploadUrl, path string, size, byteLimit int64, headers, formMap map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "files[]", size, byteLimit, formMap, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj UploadResp
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	file := obj[0]
	fileErr := file.Error
	if fileErr != nil {
		return "", errors.New("Bad response: " + fileErr.(string))
	}
	if int64(file.Size) != size {
		return "", errors.New("Byte count mismatch.")
	}
	return file.URL, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "10GB")
	if err != nil {
		return "", err
	}
	meta, err := extractMeta()
	if err != nil {
		fmt.Println("Failed to extract meta.")
		return "", err
	}
	headers := map[string]string{
		"Referer": referer,
	}
	formMap := map[string]string{
		"_sessionid":   meta.SessionID,
		"cTracker":     meta.Tracker,
		"maxChunkSize": "100000000",
		"folderId":     "-1",
		"uploadSource": "file_manager",
	}
	fileUrl, err := upload(meta.UploadURL, path, size, args.ByteLimit, headers, formMap)
	return fileUrl, err
}
