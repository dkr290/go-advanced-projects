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

type V2JsonRequest struct {
	Database string            `json:"database"`
	Key      string            `json:"key"`
	Value    map[string]string `json:"value"` // Change to a map for key-value pairs
}
type KvJsonV2 struct {
	Key   string            `json:"key"`
	Value map[string]string `json:"value"`
}
