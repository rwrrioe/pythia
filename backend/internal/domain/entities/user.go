package entities

type User struct {
	Email       string
	Name        string
	Level       string
	Lang        string
	WordsPerDay int
}

type UserStats struct {
	WordsLearned int
	AvgAccuracy  int
	TotalSession int
	Preferences  UserLearningPreferences
}

type UserLearningPreferences struct {
	Lang      string
	Level     string
	DailyGoal int
}
