package wetransfer

import "time"

type Transport struct{}

type WriteCounter struct {
	Total      int64
	TotalStr   string
	Uploaded   int64
	Percentage int
	StartTime  int64
}

type File struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Retries   int    `json:"retries"`
	Size      int64  `json:"size"`
	ItemType  string `json:"item_type"`
	ChunkSize int    `json:"chunk_size"`
}

type InitPost struct {
	Message      string `json:"message"`
	DisplayName  string `json:"display_name"`
	UILanguage   string `json:"ui_language"`
	DomainUserID string `json:"domain_user_id"`
	Files        []File `json:"files"`
}

type FileMetaResp struct {
	ID                  string      `json:"id"`
	State               string      `json:"state"`
	TransferType        int         `json:"transfer_type"`
	ShortenedURL        interface{} `json:"shortened_url"`
	RecommendedFilename string      `json:"recommended_filename"`
	ExpiresAt           time.Time   `json:"expires_at"`
	PasswordProtected   bool        `json:"password_protected"`
	UploadedAt          interface{} `json:"uploaded_at"`
	ExpiryInSeconds     int         `json:"expiry_in_seconds"`
	Size                interface{} `json:"size"`
	DeletedAt           interface{} `json:"deleted_at"`
	AccountID           interface{} `json:"account_id"`
	SecurityHash        string      `json:"security_hash"`
	From                interface{} `json:"from"`
	Creator             struct {
		Auth0UserID interface{} `json:"auth0_user_id"`
		Email       interface{} `json:"email"`
	} `json:"creator"`
	Message           string        `json:"message"`
	NumberOfDownloads int           `json:"number_of_downloads"`
	DisplayName       string        `json:"display_name"`
	Files             []File        `json:"files"`
	Recipients        []interface{} `json:"recipients"`
}

type Chunk struct {
	ChunkNumber int    `json:"chunk_number"`
	ChunkSize   int    `json:"chunk_size"`
	ChunkCrc    uint32 `json:"chunk_crc"`
	Retries     int    `json:"retries"`
}

type FilePut struct {
	URL string `json:"url"`
}

type FinaliseMppPost struct {
	ChunkCount int `json:"chunk_count"`
}
