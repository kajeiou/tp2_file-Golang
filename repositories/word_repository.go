package repositories

import (
	"tp2/interfaces"
)

var (
	DB interfaces.WordRepository
)

// AddWordToDB ajoute un nouveau mot à la base de données.
func AddWordToDB(word, definition string) error {
	return DB.AddWordToDB(word, definition)
}

// DeleteWordFromDB supprime un mot de la base de données.
func DeleteWordFromDB(word string) error {
	return DB.DeleteWordFromDB(word)
}

// ListWordsFromDB récupère tous les mots depuis la base de données.
func ListWordsFromDB() ([]interfaces.Word, error) {
	return DB.ListWordsFromDB()
}
