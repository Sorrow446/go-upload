package uguu

type Upload struct {
	Success bool `json:"success"`
	Files   []struct {
		Hash string `json:"hash"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Size string  `json:"size"`
	} `json:"files"`
}
