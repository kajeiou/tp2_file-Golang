package interfaces

type Word struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type WordRepository interface {
	InitializeDB(dbPath string) error
	CloseDB()
	ListWordsFromDB() ([]Word, error)
	AddWordToDB(word, definition string) error
	DeleteWordFromDB(word string) error
	UpdateWordInDB(word, newDefinition string) error
	GetWordFromDB(word string) (Word, error)
}
