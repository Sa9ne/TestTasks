package models

type Files struct {
	Url      string `json:"url"`
	FileName string `json:"filename"`
	Status   string `json:"status"` // Processing, Error, Done
}

type Tasks struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Links  []string `json:"links"`
	Files  []Files  `json:"files"`
	Status string   `json:"status"` // Processing, Done
}
