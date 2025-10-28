package entities

type UnknownWord struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Lang        string
}
