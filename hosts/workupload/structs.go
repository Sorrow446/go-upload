package workupload

type GetServerResp struct {
	Success bool `json:"success"`
	Data    struct {
		Server string `json:"server"`
	} `json:"data"`
}

type UploadResp struct {
	Files []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		Size int64  `json:"size"`
		Time struct {
			Date         string `json:"date"`
			TimezoneType int    `json:"timezone_type"`
			Timezone     string `json:"timezone"`
		} `json:"time"`
		Type         string `json:"type"`
		Downloads    int    `json:"downloads"`
		Permission   int    `json:"permission"`
		Expiration   bool   `json:"expiration"`
		Password     bool   `json:"password"`
		MaxDownloads bool   `json:"maxDownloads"`
		Comment      bool   `json:"comment"`
	} `json:"files"`
}
