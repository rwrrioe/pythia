package requests

type CreateSession struct {
	Duration   int `json:"durating"`
	WordsCount int `json:"words_count"`
	LangId     int `json:"lang_id"`
}

type SummarizeSession struct {
	Accuracy float64 `json:"accuracy"`
}
