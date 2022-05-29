package mixdrop

type UploadResp struct {
	File struct {
		Ref string `json:"ref"`
	} `json:"file"`
}
