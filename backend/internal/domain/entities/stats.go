package entities

type Dashboard struct {
	Streak         int       `json:"streak"`
	WordsLearned   int       `json:"words_learned"`
	Accuracy       int       `json:"accuracy"`
	LatestSessions []Session `json:"latest_sessions"`
}
