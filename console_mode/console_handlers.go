package console_mode

import (
	"bufio"
	"fmt"
	"strings"
	"tp2/dictionary"
)

func ActionAddAsync(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Entrez le nouveau mot : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	_, err := d.Get(word)
	if err == nil {
		fmt.Printf("Le mot '%s' existe déjà dans le dico.\n", word)
		return
	}

	fmt.Print("Entrez la nouvelle définition : ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.AddAsync(word, definition)

	fmt.Printf("Le mot '%s' avec la définition '%s' a été ajouté.\n", word, definition)
}

func ActionDefineAsync(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Entrez le mot : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	existingWord, err := d.Get(word)
	if err != nil {
		fmt.Printf("Le mot '%s' n'existe pas dans le dico.\n", word)
		return
	}

	fmt.Print("Entrez la nouvelle définition : ")
	newDefinition, _ := reader.ReadString('\n')
	newDefinition = strings.TrimSpace(newDefinition)

	d.EditAsync(existingWord.Word, newDefinition)

	fmt.Printf("La définition pour le mot '%s' a été mise à jour.\n", word)
}

func ActionRemoveAsync(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Écrivez le mot à supprimer : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.RemoveAsync(word)
}

func ActionList(d *dictionary.Dictionary) {
	wordsList := d.List()
	if len(wordsList) == 0 {
		fmt.Println("Aucun mot dans le dico.")
	} else {
		fmt.Println("Liste des mots du dico:")
		for _, word := range wordsList {
			fmt.Println(word.String())
		}
	}
}
