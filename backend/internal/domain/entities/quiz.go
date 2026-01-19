package entities

type QuizQuestion struct {
	Answer   string   `json:"answer"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}
