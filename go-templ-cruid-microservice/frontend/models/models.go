package models

type JsonTask struct {
	Id   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}
