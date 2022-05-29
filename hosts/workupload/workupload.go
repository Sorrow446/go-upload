package workupload

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/utils"
	"net/url"
)

const (
	referer         = "https://workupload.com/"
	getServerUrl    = referer + "api/file/getUploadServer"
	finaliseUrl     = referer + "generateLink"
	fileBagRegexStr = `name="filebag" value="([^"]+)"`
)

// Token is also in the html. File bag is epoch big endian + ?.
func getTokenAndBag() (string, string, error) {
	var token string
	respBody, err := utils.GetHtml(referer)
	if err != nil {
		return "", "", err
	}
	u, err := url.Parse(referer)
	if err != nil {
		return "", "", err
	}
	for _, c := range utils.GetCookies(u) {
		if c.Name == "token" {
			token = c.Value
			break
		}
	}
	if token == "" {
		return "", "", errors.New("The server didn't set the token cookie.")
	}
	match := utils.FindStringSubmatch(respBody, fileBagRegexStr)
	if match == nil {
		return "", "", errors.New("No regex match for file bag.")
	}
	return token, match[1], nil
}

func getServer() (string, error) {
	headers := map[string]string{
		"Referer": referer,
	}
	respBody, err := utils.DoGet(getServerUrl, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()
	var obj GetServerResp
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	if !obj.Success {
		return "", errors.New("Bad response.")
	}
	return obj.Data.Server, nil
}

func upload(uploadUrl, path string, size, byteLimit int64, headers, formMap map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "files[]", size, byteLimit, formMap, nil, headers)
	if err != nil {
		return "", err
	}
	if respBody != nil {
		defer respBody.Close()
	}
	var obj UploadResp
	err = json.NewDecoder(respBody).Decode(&obj)
	if err != nil {
		return "", err
	}
	if obj.Files[0].Size != size {
		return "", errors.New("Byte count mismatch.")
	}
	return referer + "file/" + obj.Files[0].Key, nil
}

func finalise(headers, postMap map[string]string) error {
	postMap["email"] = ""
	postMap["emailText"] = ""
	postMap["g-recaptcha-response"] = ""
	postMap["password"] = ""
	postMap["maxDownloads"] = ""
	postMap["storagetime"] = ""
	respBody, err := utils.DoFormPost(finaliseUrl, postMap, headers)
	if err != nil {
		return err
	}
	respBody.Close()
	return nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "2GB")
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Referer": referer,
	}
	token, fileBag, err := getTokenAndBag()
	if err != nil {
		fmt.Println("Failed to get token and/or file bag.")
		return "", err
	}
	uploadUrl, err := getServer()
	if err != nil {
		fmt.Println("Failed to get upload server URL.")
		return "", err
	}
	formMap := map[string]string{
		"token":   token,
		"filebag": fileBag,
	}
	fileUrl, err := upload(uploadUrl, path, size, args.ByteLimit, headers, formMap)
	if err != nil {
		return "", err
	}
	err = finalise(headers, formMap)
	if err != nil {
		fmt.Println("Failed to finalise upload.")
		return "", err
	}
	return fileUrl, err
}
