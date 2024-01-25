package repositories

import (
	"tp2/interfaces"
)

var (
	DB interfaces.WordRepository
)

func AddWordToDB(word, definition string) error {
	return DB.AddWordToDB(word, definition)
}

func DeleteWordFromDB(word string) error {
	return DB.DeleteWordFromDB(word)
}

func ListWordsFromDB() ([]interfaces.Word, error) {
	return DB.ListWordsFromDB()
}
func GetWordFromDB(word string) (interfaces.Word, error) {
	return DB.GetWordFromDB(word)
}
