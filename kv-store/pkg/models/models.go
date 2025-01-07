package models

type JsonReqiest struct {
	Database string `json:"database"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}
