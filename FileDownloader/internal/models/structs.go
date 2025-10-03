package models

type Tasks struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Links  []string `json:"links"`
	Status string   `json:"status"` // Processing, Done
}
