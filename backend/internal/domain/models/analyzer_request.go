package models

type AnalyzeRequest struct {
	Text     []string `json:"text"`
	Level    string   `json:"level"`
	Lang     string   `json:"lang"`
	Durating string   `json:"durating"`
	Book     string   `json:"book"`
	TaskID   string   `json:"taskID"`
}
