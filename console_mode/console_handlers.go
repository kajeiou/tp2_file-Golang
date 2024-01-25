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

	fmt.Print("Entrez la nouvelle définition : ")
	newDefinition, _ := reader.ReadString('\n')
	newDefinition = strings.TrimSpace(newDefinition)

	err := d.EditAsync(word, newDefinition)
	if err != nil {
		fmt.Printf("Erreur lors de la mise à jour du mot '%s' : %v\n", word, err)
		return
	}

	fmt.Printf("La définition pour le mot '%s' a été mise à jour.\n", word)
}

func ActionRemoveAsync(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Écrivez le mot à supprimer : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	err := d.RemoveAsync(word)
	if err != nil {
		fmt.Printf("Erreur lors de la suppression du mot '%s': %v\n", word, err)
	} else {
		fmt.Printf("Le mot '%s' a été supprimé avec succès.\n", word)
	}
}

func ActionList(d *dictionary.Dictionary) {
	wordsList, err := d.List()
	if err != nil {
		fmt.Printf("Erreur lors de la récupération de la liste des mots : %s\n", err.Error())
		return
	}

	if len(wordsList) == 0 {
		fmt.Println("Aucun mot dans le dico.")
	} else {
		fmt.Println("Liste des mots du dico:")
		for _, word := range wordsList {
			fmt.Println(word.String())
		}
	}
}
