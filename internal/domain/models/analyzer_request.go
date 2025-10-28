package models

type AnalyzeRequest struct {
	Text     []byte
	Level    string
	Lang     string
	Durating string
	Book     string
}
