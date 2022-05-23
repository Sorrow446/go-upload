package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
	"(KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"

var (
	jar, _ = cookiejar.New(nil)
	client = &http.Client{Transport: &Transport{}, Jar: jar}
)

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(
		"User-Agent", userAgent,
	)
	return http.DefaultTransport.RoundTrip(req)
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	var speed int64 = 0
	n := len(p)
	wc.Uploaded += int64(n)
	percentage := float64(wc.Uploaded) / float64(wc.Total) * float64(100)
	wc.Percentage = int(percentage)
	// Because of form data size.
	if wc.Uploaded > wc.Total {
		wc.Uploaded = wc.Total
		wc.Percentage = 100
	}
	toDivideBy := time.Now().UnixMilli() - wc.StartTime
	if toDivideBy != 0 {
		speed = int64(wc.Uploaded) / toDivideBy * 1000
	}
	fmt.Printf("\r%d%% @ %s/s, %s/%s ", wc.Percentage, humanize.Bytes(uint64(speed)),
		humanize.Bytes(uint64(wc.Uploaded)), wc.TotalStr)
	return n, nil
}

func CheckSize(path, sizeString string) (int64, error) {
	resolved, ok := sizeMap[sizeString]
	if !ok {
		return -1, errors.New("Invalid size limit.")
	}
	stat, err := os.Stat(path)
	if err != nil {
		return -1, err
	}
	size := stat.Size()
	if size == 0 {
		return -1, errors.New("File is empty.")
	} else if size > resolved {
		errString := fmt.Sprintf("File exceeds %s size limit.", sizeString)
		return -1, errors.New(errString)
	}
	return size, nil
}

func makeFormPart(m *multipart.Writer, fileField, filename string) (io.Writer, error) {
	mimeType := guessMimeType(filename)
	header := make(textproto.MIMEHeader)
	disposition := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fileField, filename)
	header.Set("Content-Disposition", disposition)
	header.Set("Content-Type", mimeType)
	return m.CreatePart(header)
}

func MultipartUpload(uploadUrl, path, fileField string, size, byteLimit int64, formMap, params, headers map[string]string) (io.ReadCloser, error) {
	filename := filepath.Base(path)
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	f, err := os.Open(path)
	if err != nil {
		w.Close()
		m.Close()
		return nil, err
	}
	defer f.Close()
	counter := &WriteCounter{
		Total:     size,
		TotalStr:  humanize.Bytes(uint64(size)),
		StartTime: time.Now().UnixMilli(),
	}
	// Implement and get err channel working. Seems to hang. Implement content len.
	go func() {
		defer w.Close()
		defer m.Close()
		for k, v := range formMap {
			formField, err := m.CreateFormField(k)
			if err != nil {
				return
			}
			_, err = formField.Write([]byte(v))
			if err != nil {
				return
			}
		}
		part, err := makeFormPart(m, fileField, filename)
		if err != nil {
			return
		}
		if byteLimit == -1000000 {
			_, err = io.Copy(part, f)
			if err != nil {
				return
			}
		} else {
			for range time.Tick(time.Second * 1) {
				_, err = io.CopyN(part, f, byteLimit)
				if errors.Is(err, io.EOF) {
					err = nil
					break
				}
				if err != nil {
					return
				}
			}
		}
	}()
	req, err := http.NewRequest(http.MethodPost, uploadUrl, io.TeeReader(r, counter))
	if err != nil {
		return nil, err
	}
	if headers != nil {
		setHeaders(req, headers)
	}
	req.Header.Add("Content-Type", m.FormDataContentType())
	if params != nil {
		setParams(req, params)
	}
	defer fmt.Println("")
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !(do.StatusCode == http.StatusOK || do.StatusCode == http.StatusCreated) {
		do.Body.Close()
		return nil, errors.New(do.Status)
	}
	return do.Body, nil
}

// func PutUpload(uploadUrl, path string, size int64, params, headers map[string]string) (io.ReadCloser, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	counter := &WriteCounter{Total: size, TotalStr: humanize.Bytes(uint64(size))}
// 	req, err := http.NewRequest(http.MethodPut, uploadUrl, io.TeeReader(f, counter))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if params != nil {
// 		setParams(req, params)
// 	}
// 	if headers != nil {
// 		setHeaders(req, headers)
// 	}
// 	mimeType := guessMimeType(path)
// 	req.Header.Add("Content-Type", mimeType)
// 	do, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !(do.StatusCode == http.StatusOK || do.StatusCode == http.StatusCreated) {
// 		do.Body.Close()
// 		return nil, errors.New(do.Status)
// 	}
// 	return do.Body, nil
// }

// Do headers and params.
func GetHtml(url string) (string, error) {
	reqBody, err := DoGet(url, nil, nil)
	if err != nil {
		return "", err
	}
	defer reqBody.Close()
	bodyBytes, err := io.ReadAll(reqBody)
	return string(bodyBytes), err
}

func FindStringSubmatch(text, regexString string) []string {
	regex := regexp.MustCompile(regexString)
	match := regex.FindStringSubmatch(text)
	return match
}

func FindStringSubmatches(text, regexString string) [][]string {
	regex := regexp.MustCompile(regexString)
	matches := regex.FindAllStringSubmatch(text, -1)
	return matches
}

func FindHtmlSubmatch(url, regexString string) ([]string, error) {
	html, err := GetHtml(url)
	if err != nil {
		return nil, err
	}
	match := FindStringSubmatch(html, regexString)
	return match, nil
}

func FindHtmlSubmatches(_url, regexString string) ([][]string, error) {
	html, err := GetHtml(_url)
	if err != nil {
		return nil, err
	}
	matches := FindStringSubmatches(html, regexString)
	return matches, nil
}

func setParams(req *http.Request, params map[string]string) {
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	req.URL.RawQuery = query.Encode()
}

func setHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}

func DoGet(_url string, params, headers map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		setHeaders(req, headers)
	}
	if params != nil {
		setParams(req, params)
	}
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if do.StatusCode != http.StatusOK {
		do.Body.Close()
		return nil, errors.New(do.Status)
	}
	return do.Body, nil
}

func makeJsonReq(_url string, jsonMap map[string]interface{}) (*http.Request, error) {
	m, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(m))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	return req, nil
}

func DoPost(_url string, params, headers map[string]string, jsonMap map[string]interface{}) (io.ReadCloser, error) {
	var (
		err error
		req *http.Request
	)
	if jsonMap != nil {
		req, err = makeJsonReq(_url, jsonMap)
	} else {
		req, err = http.NewRequest(http.MethodPost, _url, nil)
	}
	if err != nil {
		return nil, err
	}
	if headers != nil {
		setHeaders(req, headers)
	}
	if params != nil {
		setParams(req, params)
	}
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if do.StatusCode != http.StatusOK {
		do.Body.Close()
		return nil, errors.New(do.Status)
	}
	return do.Body, nil
}

// Unknown extensions and extensions with multiple periods in will be treated as octet.
func guessMimeType(path string) string {
	octetMime := "application/octet-stream"
	lastIndex := strings.LastIndex(path, ".")
	if lastIndex == -1 {
		return octetMime
	}
	extension := path[lastIndex:]
	resolved, ok := mimeMap[extension]
	if !ok {
		return octetMime
	}
	return resolved
}
