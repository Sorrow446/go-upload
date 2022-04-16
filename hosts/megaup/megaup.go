package megaup

import (
	"encoding/json"
	"errors"
	"main/utils"
	"strconv"
)

const (
	referer        = "https://megaup.net/"
	urlRegexString = `https://f\d{1,3}.megaup.net/core/page/ajax/file_upload_handler.ajax.php\?` +
		`r=megaup.net&p=https&csaKey1=[a-z\d]{64}&csaKey2=[a-z\d]{64}`
	sessTrackRegexString = `_sessionid: '([a-z\d]{26})'. cTracker: '([a-z\d]{32})'`
)

func getUploadUrl() (string, string, string, error) {

	html, err := utils.GetHtml(referer)
	if err != nil {
		return "", "", "", err
	}
	match := utils.FindStringSubmatch(html, urlRegexString)
	if err != nil {
		return "", "", "", err
	}
	if match == nil {
		return "", "", "", errors.New("No regex match.")
	}

	match = utils.FindStringSubmatch(html, sessTrackRegexString)
	if err != nil {
		return "", "", "", err
	}
	if match == nil {
		return "", "", "", errors.New("No regex match.")
	}
	return match[0], match[1], match[2], nil
}

func upload(uploadUrl, path string, size, byteLimit int64, formMap, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "files[]", size, byteLimit, formMap, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj Upload
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	returnedSize, err := strconv.ParseInt(obj[0].Size, 10, 64)
	if err != nil {
		return "", err
	}
	if obj[0].Error != nil {
		return "", errors.New("Bad response.")
	} else if returnedSize != size {
		return "", errors.New("Byte count mismatch.")
	}
	return obj[0].URL, nil
}

func Run(args *utils.Args, path string) (string, error) {
	uploadUrl, sessionId, tracker, err := getUploadUrl()
	if err != nil {
		return "", err
	}
	size, err := utils.CheckSize(path, "5GB")
	if err != nil {
		return "", err
	}
	formMap := map[string]string{
		"_sessionid":   sessionId,
		"cTracker":     tracker,
		"folderId":     "",
		"maxChunkSize": "100000000",
	}
	headers := map[string]string{
		"Referer": referer,
	}
	fileUrl, err := upload(uploadUrl, path, size, args.ByteLimit, formMap, headers)
	return fileUrl, err
}
