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
	addCh      chan Word // Canal pour ajouter un mot de manière asynchrone
	editCh     chan Word
	removeCh   chan string   // Canal pour supprimer un mot de manière asynchrone
	mu         sync.Mutex    // Mutex pour éviter les problèmes de concurrence
	responseCh chan struct{} // Canal pour signaler la fin d'une opération asynchrone
}

func (w Word) String() string {

	return w.Word + ": " + w.Definition
}

func New(filename string) *Dictionary {
	d := &Dictionary{
		filename:   filename,
		words:      make([]Word, 0),
		addCh:      make(chan Word),
		editCh:     make(chan Word),
		removeCh:   make(chan string),
		responseCh: make(chan struct{}),
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

// AddAsync ajoute de manière asynchrone un nouveau mot.
func (d *Dictionary) AddAsync(word string, definition string) {

	d.mu.Lock()
	// Verrouille le mutex pour assurer un accès exclusif aux données du dico.
	//  Si un autre processus ou une autre goroutine tente d'appeler AddAsync ou toute autre fonction qui modifie les données partagées,
	// elle devra attendre que le mutex soit déverrouillé avant de pouvoir procéder
	defer d.mu.Unlock()

	// defer signifie que l'instruction d.mu.Unlock() sera exécutée lorsque AddAsync prend fin
	// Cela garantit que le mutex est déverrouillé, même si une panique survient (une panique est une situation exceptionnelle qui peut se produire en cas d'erreur grave).
	// Si le mutex n'était pas déverrouillé en cas de panique, cela pourrait entraîner un verrouillage permanent du mutex, rendant l'ensemble du programme inutilisable.

	newWord := Word{Word: word, Definition: definition}
	d.words = append(d.words, newWord)

	go func() {
		d.responseCh <- struct{}{} // Signale la fin de l'opération
	}()
}
func (d *Dictionary) EditAsync(word string, newDefinition string) {

	d.mu.Lock()
	// Verrouille le mutex pour assurer un accès exclusif aux données du dico.
	//  Si un autre processus ou une autre goroutine tente d'appeler AddAsync ou toute autre fonction qui modifie les données partagées,
	// elle devra attendre que le mutex soit déverrouillé avant de pouvoir procéder
	defer d.mu.Unlock()

	// defer signifie que l'instruction d.mu.Unlock() sera exécutée lorsque AddAsync prend fin
	// Cela garantit que le mutex est déverrouillé, même si une panique survient (une panique est une situation exceptionnelle qui peut se produire en cas d'erreur grave).
	// Si le mutex n'était pas déverrouillé en cas de panique, cela pourrait entraîner un verrouillage permanent du mutex, rendant l'ensemble du programme inutilisable.

	for i, w := range d.words {
		if w.Word == word {
			d.words[i].Definition = newDefinition
			break
		}
	}

	go func() {
		d.responseCh <- struct{}{} // Signale la fin de l'opération
	}()
}

// RemoveAsync supprime de manière asynchrone un mot
func (d *Dictionary) RemoveAsync(word string) {

	d.mu.Lock()
	// Verrouille le mutex pour assurer un accès exclusif aux données du dico.
	//  Si un autre processus ou une autre goroutine tente d'appeler AddAsync ou toute autre fonction qui modifie les données partagées,
	// elle devra attendre que le mutex soit déverrouillé avant de pouvoir procéder
	defer d.mu.Unlock()

	// defer signifie que l'instruction d.mu.Unlock() sera exécutée lorsque AddAsync prend fin
	// Cela garantit que le mutex est déverrouillé, même si une panique survient (une panique est une situation exceptionnelle qui peut se produire en cas d'erreur grave).
	// Si le mutex n'était pas déverrouillé en cas de panique, cela pourrait entraîner un verrouillage permanent du mutex, rendant l'ensemble du programme inutilisable.

	var updatedWords []Word
	for _, w := range d.words {
		if w.Word != word {
			updatedWords = append(updatedWords, w)
		}
	}
	d.words = updatedWords
	d.responseCh <- struct{}{}
}

// Get retourne un mot du dico en fonction du mot fourni.
func (d *Dictionary) Get(word string) (Word, error) {
	for _, w := range d.words {
		if w.Word == word {
			return w, nil // Retourne le mot si trouvé
		}
	}

	return Word{}, errors.New("Le mot " + word + " n'a pas été trouvé dans le dico")
}

// List retourne la liste complète des mots dans le dico.
func (d *Dictionary) List() []Word {
	wordsList := make([]string, 0)
	for _, w := range d.words {
		wordsList = append(wordsList, w.Definition)
	}
	return d.words
}

// chargerFichier charge le contenu du fichier CSV dans le dictionnaire
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
