package requests

type AnalyzeRequest struct {
	Level    string `json:"level"`
	Durating string `json:"durating"`
	Lang     string `json:"lang"`
}
