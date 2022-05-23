// Why not just one simple form post, WeTransfer devs?

package wetransfer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"main/utils"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	defaultChunkSize = 5242880
	referer          = "https://wetransfer.com/"
	userAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"
	apiBase      = referer + "api/v4/transfers/"
	apiLinkUrl   = apiBase + "link"
	csrfTokenStr = `name="csrf-token" content="([^"]+)"`
)

var (
	jar, _ = cookiejar.New(nil)
	client = &http.Client{Transport: &Transport{}, Jar: jar}
)

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(
		"User-Agent", userAgent,
	)
	req.Header.Add(
		"Referer", referer,
	)
	req.Header.Add(
		"Origin", referer,
	)
	return http.DefaultTransport.RoundTrip(req)
}

func (wc *WriteCounter) printProgress(n int) (int, error) {
	var speed int64 = 0
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

func getCsrfToken() (string, error) {
	req, err := client.Get(referer)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return "", errors.New(req.Status)
	}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	regex := regexp.MustCompile(csrfTokenStr)
	match := regex.FindStringSubmatch(string(bodyBytes))
	if match == nil {
		return "", errors.New("No regex match.")
	}
	return match[1], nil
}

func addHeaders(req *http.Request, csrfToken string) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-CSRF-Token", csrfToken)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
}

func initiate(csrfToken string, _file *File) (string, error) {
	_file.ItemType = "file"
	postData := InitPost{
		DisplayName: _file.Name,
		Message:     "",
		UILanguage:  "en",
		Files:       []File{*_file},
	}
	m, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, apiLinkUrl, bytes.NewBuffer(m))
	if err != nil {
		return "", err
	}
	addHeaders(req, csrfToken)
	do, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return "", errors.New(do.Status)
	}
	var obj FileMetaResp
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return "", err
	}
	if obj.State != "processing" {
		return "", errors.New("Invalid state.")
	}
	return obj.ID, nil
}

func getFileId(csrfToken, transferId string, _file *File) (string, error) {
	_url := apiBase + transferId + "/files"
	m, err := json.Marshal(&_file)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(m))
	if err != nil {
		return "", err
	}
	addHeaders(req, csrfToken)
	do, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return "", errors.New(do.Status)
	}
	var obj File
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return "", err
	}
	return obj.ID, nil
}

func uploadChunks(csrfToken, transferId, fileId, filePath string, size int64) (int, error) {
	_url := apiBase + transferId + "/files/" + fileId + "/part-put-url"
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	chunk := make([]byte, defaultChunkSize)
	chunkNum := 0
	totalChunkSize := 0
	counter := &WriteCounter{
		Total:     size,
		TotalStr:  humanize.Bytes(uint64(size)),
		StartTime: time.Now().UnixMilli(),
	}
	for {
		chunkSize, err := reader.Read(chunk)
		if chunkNum == 0 {
			defer fmt.Println("")
			counter.printProgress(0)
		}
		totalChunkSize += chunkSize
		if errors.Is(err, io.EOF) || chunkSize == 0 {
			break
		} else if err != nil {
			return -1, err
		}
		chunkNum++
		postData := Chunk{
			ChunkSize:   chunkSize,
			ChunkNumber: chunkNum,
			ChunkCrc:    crc32.ChecksumIEEE(chunk[:chunkSize]),
			Retries:     0,
		}
		m, err := json.Marshal(postData)
		if err != nil {
			return -1, err
		}
		req, err := http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(m))
		if err != nil {
			return -1, err
		}
		addHeaders(req, csrfToken)
		do, err := client.Do(req)
		if err != nil {
			return -1, err
		}
		if do.StatusCode != http.StatusOK {
			do.Body.Close()
			return -1, errors.New(do.Status)
		}
		var obj FilePut
		err = json.NewDecoder(do.Body).Decode(&obj)
		do.Body.Close()
		if err != nil {
			return -1, err
		}
		// Using TeeReader here causes 501s.
		req2, err := http.NewRequest(http.MethodPut, obj.URL, bytes.NewBuffer(chunk[:chunkSize]))
		if err != nil {
			return -1, err
		}
		req2.Header.Set("Content-Type", "binary/octet-stream")
		req2.Header.Set("Content-Length", strconv.Itoa(chunkSize))
		do2, err := client.Do(req2)
		if err != nil {
			return -1, err
		}
		do2.Body.Close()
		if do2.StatusCode != http.StatusOK {
			return -1, errors.New(do2.Status)
		}
		counter.printProgress(chunkSize)
	}
	return chunkNum, nil
}

func finalise(csrfToken, transferId, fileId string, chunkCount int, fileSize int64) (string, error) {
	_url := apiBase + transferId + "/files/" + fileId + "/finalize-mpp"
	postData := FinaliseMppPost{
		ChunkCount: chunkCount,
	}
	m, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPut, _url, bytes.NewBuffer(m))
	if err != nil {
		return "", err
	}
	addHeaders(req, csrfToken)
	do, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return "", errors.New(do.Status)
	}
	_url2 := apiBase + "/" + transferId + "/finalize"
	req2, err := http.NewRequest(http.MethodPut, _url2, nil)
	if err != nil {
		return "", err
	}
	addHeaders(req2, csrfToken)
	do2, err := client.Do(req2)
	if err != nil {
		return "", err
	}
	defer do2.Body.Close()
	if do2.StatusCode != http.StatusOK {
		return "", errors.New(do2.Status)
	}
	var obj FileMetaResp
	err = json.NewDecoder(do2.Body).Decode(&obj)
	if err != nil {
		return "", err
	}
	if obj.Files[0].Size != fileSize {
		return "", errors.New("Byte count mismatch.")
	}
	return obj.ShortenedURL.(string), nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "2GB")
	if err != nil {
		return "", err
	}
	csrfToken, err := getCsrfToken()
	if err != nil {
		fmt.Println("Failed to get CSRF token.")
		return "", err
	}
	_file := &File{
		Name: filepath.Base(path),
		Size: size,
	}
	transferId, err := initiate(csrfToken, _file)
	if err != nil {
		fmt.Println("Failed to initiate upload.")
		return "", err
	}
	fileId, err := getFileId(csrfToken, transferId, _file)
	if err != nil {
		fmt.Println("Failed to get file ID.")
		return "", err
	}
	chunkCount, err := uploadChunks(csrfToken, transferId, fileId, path, size)
	if err != nil {
		fmt.Println("Failed to upload file chunks.")
		return "", err
	}
	fileUrl, err := finalise(csrfToken, transferId, fileId, chunkCount, size)
	if err != nil {
		fmt.Println("Failed to get finalise upload.")
		return "", err
	}
	return fileUrl, err
}
