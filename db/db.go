package db

import (
	"tp2/dictionary"
	"tp2/interfaces"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormWordRepository struct {
	DB *gorm.DB
}

func (g *GormWordRepository) InitializeDB(dbPath string) error {
	var err error
	g.DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	g.DB.AutoMigrate(&dictionary.Word{})
	return nil
}

func (g *GormWordRepository) CloseDB() {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

func (g *GormWordRepository) AddWordToDB(word, definition string) error {
	newWord := dictionary.Word{
		Word:       word,
		Definition: definition,
	}

	result := g.DB.Create(&newWord)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (g *GormWordRepository) DeleteWordFromDB(word string) error {
	result := g.DB.Where("word = ?", word).Unscoped().Delete(&dictionary.Word{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (g *GormWordRepository) ListWordsFromDB() ([]interfaces.Word, error) {
	var words []dictionary.Word
	result := g.DB.Find(&words)
	if result.Error != nil {
		return nil, result.Error
	}
	var interfaceWords []interfaces.Word
	for _, w := range words {
		interfaceWords = append(interfaceWords, interfaces.Word{
			Word:       w.Word,
			Definition: w.Definition,
		})
	}

	return interfaceWords, nil
}

func (g *GormWordRepository) UpdateWordInDB(word, newDefinition string) error {
	var existingWord dictionary.Word
	result := g.DB.Where("word = ?", word).First(&existingWord)
	if result.Error != nil {
		return result.Error
	}

	existingWord.Definition = newDefinition

	result = g.DB.Save(&existingWord)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *GormWordRepository) GetWordFromDB(word string) (interfaces.Word, error) {
	var existingWord dictionary.Word
	result := g.DB.Where("word = ?", word).First(&existingWord)
	if result.Error != nil {
		return interfaces.Word{}, result.Error
	}

	return interfaces.Word{
		Word:       existingWord.Word,
		Definition: existingWord.Definition,
	}, nil
}
