package onefichier

type GetServer struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Finalize struct {
	Incoming int `json:"incoming"`
	Links    []struct {
		Download  string `json:"download"`
		Filename  string `json:"filename"`
		Remove    string `json:"remove"`
		Size      string `json:"size"`
		Whirlpool string `json:"whirlpool"`
	} `json:"links"`
}
