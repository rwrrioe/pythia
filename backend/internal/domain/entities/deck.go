package entities

type Deck struct {
	Id         int
	SessionId  int64
	Flashcards []FlashCard
}
