package dictionary

import (
	"encoding/csv"
	"errors"
	"os"
	"sync"
	"tp2/interfaces"

	"gorm.io/gorm"
)

type Word struct {
	gorm.Model `gorm:"soft_delete:false"`
	Word       string `gorm:"unique;not null"`
	Definition string `gorm:"not null"`
}

type Dictionary struct {
	filename   string
	words      []Word
	addCh      chan Word // Canal pour ajouter un mot de manière asynchrone
	editCh     chan Word
	removeCh   chan string               // Canal pour supprimer un mot de manière asynchrone
	mu         sync.Mutex                // Mutex pour éviter les problèmes de concurrence
	responseCh chan struct{}             // Canal pour signaler la fin d'une opération asynchrone
	wordRepo   interfaces.WordRepository // Ajouter le champ wordRepo à la structure Dictionary
}

func (w Word) String() string {

	return w.Word + ": " + w.Definition
}

func New(filename string, wordRepository interfaces.WordRepository) *Dictionary {
	d := &Dictionary{
		filename:   filename,
		words:      make([]Word, 0),
		addCh:      make(chan Word),
		editCh:     make(chan Word),
		removeCh:   make(chan string),
		responseCh: make(chan struct{}),
		wordRepo:   wordRepository,
	}
	go d.processChannels() // Lance la gestion asynchrone des canaux
	d.chargerFichier()     // Charge le dico depuis le fichier
	return d
}

// processChannels gère de manière asynchrone les opérations sur les canaux.
func (d *Dictionary) processChannels() {
	for {
		select {
		case word := <-d.addCh:
			d.AddAsync(word.Word, word.Definition) // Ajoute de manière asynchrone un nouveau mot
			<-d.responseCh                         // Attend la fin de l'opération
		case word := <-d.editCh:
			d.EditAsync(word.Word, word.Definition) // Modifie de manière asynchrone un nouveau mot
			<-d.responseCh                          // Attend la fin de l'opération
		case word := <-d.removeCh:
			d.RemoveAsync(word) // Supprime de manière asynchrone un mot
			<-d.responseCh      // Attend la fin de l'opération
		case <-d.responseCh:
			d.enregistrerFichier() // Enregistre le dico dans le fichier après une opération
		}
	}
}

func (d *Dictionary) AddAsync(word string, definition string) error {
	// Mutex pour synchroniser l'accès à d.mu
	d.mu.Lock()
	defer d.mu.Unlock()

	// Utilise la méthode AddWordToDB du repository wordRepo pour ajouter le mot à la base de données
	if err := d.wordRepo.AddWordToDB(word, definition); err != nil {
		// En cas d'erreur lors de l'ajout à la base de données, retourne l'erreur
		return err
	}

	// Pas d'erreur, signale la fin de l'opération
	d.responseCh <- struct{}{}
	return nil
}

// GetResponseChannel renvoie le canal de réponse du dictionnaire.
func (d *Dictionary) ResponseChannel() <-chan struct{} {
	return d.responseCh
}

func (d *Dictionary) EditAsync(word string, newDefinition string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Utilise la méthode GetWordFromDB du repository pour obtenir le mot depuis la base de données
	existingWord, err := d.wordRepo.GetWordFromDB(word)
	if err != nil {
		// Gère l'erreur si le mot n'est pas trouvé dans la base de données
		return err
	}

	existingWord.Definition = newDefinition

	if err := d.wordRepo.UpdateWordInDB(existingWord.Word, existingWord.Definition); err != nil {
		return err
	}

	return nil
}

// RemoveAsync supprime de manière asynchrone un mot
func (d *Dictionary) RemoveAsync(word string) error {
	// Mutex pour synchroniser l'accès à d.mu
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.wordExists(word) {
		return errors.New("Le mot n'existe pas dans le dictionnaire")
	}

	if err := d.wordRepo.DeleteWordFromDB(word); err != nil {
		// En cas d'erreur lors de la suppression de la base de données, signale la fin de l'opération avec une erreur
		d.responseCh <- struct{}{}
		return err
	}
	d.responseCh <- struct{}{}
	return nil
}

func (d *Dictionary) wordExists(word string) bool {
	_, err := d.wordRepo.GetWordFromDB(word)
	return err == nil
}

func (d *Dictionary) List() ([]Word, error) {
	wordsFromDB, err := d.wordRepo.ListWordsFromDB()
	if err != nil {
		return nil, err
	}

	// Convertir []interfaces.Word en []Word
	words := make([]Word, len(wordsFromDB))
	for i, w := range wordsFromDB {
		words[i] = Word{Word: w.Word, Definition: w.Definition}
	}

	return words, nil
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

// enregistrerFichier enregistre le dico dans le fichier CSV
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
