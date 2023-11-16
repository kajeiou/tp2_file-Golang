package dictionary

import (
	"encoding/csv"
	"errors"
	"os"
	"sync"
)

type Word struct {
	Word       string
	Definition string
}

type Dictionary struct {
	filename   string
	words      []Word
	addCh      chan Word
	removeCh   chan string
	mu         sync.Mutex
	responseCh chan struct{}
}

func (w Word) String() string {

	return w.Word + ": " + w.Definition
}

func New(filename string) *Dictionary {
	d := &Dictionary{
		filename:   filename,
		words:      make([]Word, 0),
		addCh:      make(chan Word),
		removeCh:   make(chan string),
		responseCh: make(chan struct{}),
	}
	go d.processChannels()
	d.chargerFichier()
	return d
}
func (d *Dictionary) processChannels() {
	for {
		select {
		case word := <-d.addCh:
			d.AddAsync(word.Word, word.Definition)
			<-d.responseCh
		case word := <-d.removeCh:
			d.RemoveAsync(word)
		case <-d.responseCh:
			d.enregistrerFichier()
		}
	}
}

func (d *Dictionary) AddAsync(word string, definition string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	newWord := Word{Word: word, Definition: definition}
	d.words = append(d.words, newWord)

	go func() {
		d.responseCh <- struct{}{}
	}()
}

func (d *Dictionary) RemoveAsync(word string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var updatedWords []Word
	for _, w := range d.words {
		if w.Word != word {
			updatedWords = append(updatedWords, w)
		}
	}
	d.words = updatedWords
	d.responseCh <- struct{}{}
}

func (d *Dictionary) Get(word string) (Word, error) {
	for _, w := range d.words {
		if w.Word == word {
			return w, nil
		}
	}

	return Word{}, errors.New("Le mot " + word + " n'a pas été trouvé dans le dico")
}

func (d *Dictionary) List() []Word {
	wordsList := make([]string, 0)
	for _, w := range d.words {
		wordsList = append(wordsList, w.Definition)
	}
	return d.words
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
		if len(record) == 2 {
			word := Word{
				Word:       record[0],
				Definition: record[1],
			}
			d.words = append(d.words, word)
		}
	}

	return nil
}

func (d *Dictionary) enregistrerFichier() error {
	d.mu.Lock()
	defer d.mu.Unlock()

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
