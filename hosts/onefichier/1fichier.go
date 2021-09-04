// Needs content length header (size of file + form). MultipartUpload func doesn't support this yet.

package onefichier

import (
	"encoding/json"
	"errors"
	"main/utils"
	"strconv"
)

const apiBase = "https://www.filemail.com/api/transfer/"

func getServer(headers map[string]string) (string, string, error) {
	url := "https://api.1fichier.com/v1/upload/get_upload_server.cgi"
	headers["Content-Type"] = "application/json"
	respBody, err := utils.DoGet(url, nil, headers)
	if err != nil {
		return "", "", err
	}
	defer respBody.Close()
	var obj GetServer
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", "", err
	}
	return "https://" + obj.URL, obj.ID, nil
}

func upload(uploadUrl, path, id string, size int64, headers map[string]string) error {
	uploadUrl += "/upload.cgi"
	params := map[string]string{
		"id": id,
	}
	formMap := map[string]string{
		"send_ssl": "on",
		"domain":   "0",
		"mail":     "",
		"dpass":    "",
		"user":     "",
		"mails":    "",
		"message":  "",
		"submit":   "Send",
	}
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file[]", size, formMap, params, headers)
	if err != nil {
		return err
	}
	respBody.Close()
	return nil
}

func finalizeUpload(finalizeUrl, id string, size int64, headers map[string]string) (string, error) {
	params := map[string]string{
		"xid": id,
	}
	finalizeUrl += "/end.pl"
	headers["Content-Type"] = "application/json"
	respBody, err := utils.DoGet(finalizeUrl, params, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj Finalize
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	returnedSize, err := strconv.ParseInt(obj.Links[0].Size, 10, 64)
	if err != nil {
		return "", err
	}
	if returnedSize != size {
		return "", errors.New("Byte count mismatch.")
	}
	return obj.Links[0].Download, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "300GB")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": "https://1fichier.com/",
	}
	server, id, err := getServer(headers)
	if err != nil {
		return "", err
	}
	err = upload(server, path, id, size, headers)
	if err != nil {
		return "", err
	}
	fileUrl, err := finalizeUpload(server, id, size, headers)
	return fileUrl, err
}
