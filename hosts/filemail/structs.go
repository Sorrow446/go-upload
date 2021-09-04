package filemail

type Initialize struct {
	Transferid     string `json:"transferid"`
	Transferkey    string `json:"transferkey"`
	Transferurl    string `json:"transferurl"`
	Transferip     string `json:"transferip"`
	Udpport        int    `json:"udpport"`
	Udpthreshold   int    `json:"udpthreshold"`
	Responsestatus string `json:"responsestatus"`
}

type Finalize struct {
	Downloadurl    string `json:"downloadurl"`
	Responsestatus string `json:"responsestatus"`
}
