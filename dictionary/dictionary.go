package dictionary

import (
	"encoding/csv"
	"errors"
	"os"
)

type Word struct {
	Word       string
	Definition string
}

type Dictionary struct {
	filename string
	words    []Word
}

func New(filename string) *Dictionary {
	d := &Dictionary{
		filename: filename,
		words:    make([]Word, 0),
	}
	d.chargerFichier()
	return d
}

func (d *Dictionary) Add(word string, definition string) {
	newWord := Word{
		Word:       word,
		Definition: definition,
	}

	d.words = append(d.words, newWord)

	d.enregistrerFichier()

}

func (d *Dictionary) Get(word string) (Word, error) {
	for _, w := range d.words {
		if w.Word == word {
			return w, nil
		}
	}

	return Word{}, errors.New("Le mot " + word + " n'a pas été trouvé dans le dico")
}

func (d *Dictionary) Remove(word string) {

}

func (d *Dictionary) List() ([]string, []Word) {
	wordsList := make([]string, 0)
	for _, w := range d.words {
		wordsList = append(wordsList, w.Definition)
	}
	return wordsList, d.words
}

func (d *Dictionary) chargerFichier() error {
	file, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) == 1 {
			word := Word{Definition: record[0]}
			d.words = append(d.words, word)
		}
	}

	return nil
}

func (d *Dictionary) enregistrerFichier() error {
	file, err := os.Create(d.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, word := range d.words {
		err := writer.Write([]string{word.Word, word.Definition})
		if err != nil {
			return err
		}
	}

	return nil
}
