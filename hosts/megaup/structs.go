package megaup

type Upload []struct {
	Name              string      `json:"name"`
	Size              string      `json:"size"`
	Type              string      `json:"type"`
	Error             interface{} `json:"error"`
	URL               string      `json:"url"`
	DeleteURL         string      `json:"delete_url"`
	InfoURL           string      `json:"info_url"`
	DeleteType        string      `json:"delete_type"`
	DeleteHash        string      `json:"delete_hash"`
	Hash              string      `json:"hash"`
	StatsURL          string      `json:"stats_url"`
	ShortURL          string      `json:"short_url"`
	FileID            string      `json:"file_id"`
	UniqueHash        string      `json:"unique_hash"`
	URLHTML           string      `json:"url_html"`
	URLBbcode         string      `json:"url_bbcode"`
	SuccessResultHTML string      `json:"success_result_html"`
}
