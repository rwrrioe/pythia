package entities

type UnknownWord struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Lang        string
}

type Example struct {
	Word    string `json:"word"`
	Example string `json:"example"`
}
