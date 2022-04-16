package zippyshare

import (
	"errors"
	"io"
	"main/utils"
	"path/filepath"
)

const (
	referer     = "https://www.zippyshare.com/"
	serverRegex = `var server = \'(www\d{1,3})\';`
	urlRegex    = `onclick="this.select\(\);" value="(https://www\d{1,3}.zippyshare` +
		`.com/v/[a-zA-Z\d]{8}/file.html)`
)

func getServer() (string, error) {
	match, err := utils.FindHtmlSubmatch(referer, serverRegex)
	if err != nil {
		return "", err
	}
	if match == nil {
		return "", errors.New("No regex match.")
	}
	return match[1], nil
}

func extractUrl(html string) (string, error) {
	match := utils.FindStringSubmatch(html, urlRegex)
	if match == nil {
		return "", errors.New("No regex match.")
	}
	return match[1], nil
}

func upload(uploadUrl, path string, size, byteLimit int64, formMap, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file", size, byteLimit, formMap, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	bodyBytes, err := io.ReadAll(respBody)
	if err != nil {
		return "", err
	}
	fileUrl, err := extractUrl(string(bodyBytes))
	return fileUrl, err
}

func Run(args *utils.Args, path string) (string, error) {
	server, err := getServer()
	if err != nil {
		return "", err
	}
	uploadUrl := "https://" + server + ".zippyshare.com/upload"
	size, err := utils.CheckSize(path, "500MB")
	if err != nil {
		return "", err
	}
	formMap := map[string]string{
		"name":            filepath.Base(path),
		"zipname":         "",
		"ziphash":         "",
		"embPlayerValues": "false",
	}
	headers := map[string]string{
		"Referer": referer,
	}
	if args.Private {
		formMap["private"] = "true"
	} else {
		formMap["notprivate"] = "true"
	}
	fileUrl, err := upload(uploadUrl, path, size, args.ByteLimit, formMap, headers)
	return fileUrl, err
}
