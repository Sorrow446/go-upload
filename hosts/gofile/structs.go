package gofile

type GetServer struct {
	Status string `json:"status"`
	Data   struct {
		Server string `json:"server"`
	} `json:"data"`
}

type Upload struct {
	Status string `json:"status"`
	Data   struct {
		DownloadPage string `json:"downloadPage"`
		Code         string `json:"code"`
		ParentFolder string `json:"parentFolder"`
		FileID       string `json:"fileId"`
		FileName     string `json:"fileName"`
		Md5          string `json:"md5"`
		DirectLink   string `json:"directLink"`
		Info         string `json:"info"`
	} `json:"data"`
}
