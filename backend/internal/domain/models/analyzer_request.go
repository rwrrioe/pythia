package models

type AnalyzeRequest struct {
	Level    string `json:"level"`
	Lang     string `json:"lang"`
	Durating string `json:"durating"`
	Book     string `json:"book"`
	TaskID   string `json:"taskID"`
}
