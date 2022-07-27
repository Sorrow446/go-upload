package krakenfiles

type Upload struct {
	Files []struct {
		Name  string `json:"name"`
		Size  string `json:"size"`
		Error string `json:"error"`
		URL   string `json:"url"`
		Hash  string `json:"hash"`
	} `json:"files"`
}
