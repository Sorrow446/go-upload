package fileio

import "time"

type Upload struct {
	Success      bool      `json:"success"`
	Status       int       `json:"status"`
	ID           string    `json:"id"`
	Key          string    `json:"key"`
	Name         string    `json:"name"`
	Link         string    `json:"link"`
	Private      bool      `json:"private"`
	Expires      time.Time `json:"expires"`
	Downloads    int       `json:"downloads"`
	MaxDownloads int       `json:"maxDownloads"`
	AutoDelete   bool      `json:"autoDelete"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mimeType"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}
