package models

type JsonRequest struct {
	Database string   `json:"database"`
	Key      string   `json:"key"`
	Value    []string `json:"value"`
}
type JsonRequestGet struct {
	Database string `json:"database"`
	Key      string `json:"key"`
}

type KvJson struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}
