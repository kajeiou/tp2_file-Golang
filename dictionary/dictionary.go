package dictionary

import (
	"encoding/csv"
	"os"
)

type Word struct {
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

}

func (d *Dictionary) Get(word string) (*Word, bool) {

}

func (d *Dictionary) Remove(word string) {

}

func (d *Dictionary) List() ([]string, []Word) {

}

func (d *Dictionary) chargerFichier() error {
	file, err := os.OpenFile(d.filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
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
		err := writer.Write([]string{word.Definition})
		if err != nil {
			return err
		}
	}

	return nil
}
