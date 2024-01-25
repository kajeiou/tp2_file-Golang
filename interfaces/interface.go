package interfaces

type Word struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type WordRepository interface {
	InitializeDB() error
	CloseDB()
	ListWordsFromDB() ([]Word, error)
	AddWordToDB(word, definition string) error
	DeleteWordFromDB(word string) error
}
