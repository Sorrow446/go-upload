package racaty

type UploadResp []struct {
	FileCode   string `json:"file_code"`
	FileStatus string `json:"file_status"`
}
