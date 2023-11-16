package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tp2/dictionary"
)

func main() {
	myDictionary := dictionary.New("dictionary.csv")
	fmt.Println("Bienvenue dans le dico !")

	for {
		fmt.Println("|| MENU Dico ||")
		fmt.Println("Voir : 1")
		fmt.Println("Ajouter : 2,  Définir : 3")
		fmt.Println("Supprimer : 4, Sortir : 5")
		fmt.Println("Choisissez ...")

		reader := bufio.NewReader(os.Stdin)

		choix, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erreur de lecture de l'entrée utilisateur:", err)
			return
		}
		choix = strings.TrimSpace(choix)

		switch choix {
		case "1":
			actionList(myDictionary)
		case "2":
			actionAddAsync(myDictionary, reader)
		case "3":
			actionDefineAsync(myDictionary, reader)
		case "4":
			actionRemoveAsync(myDictionary, reader)
		case "5":
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Choix invalide. Veuillez entrer un numéro valide.")
		}
	}
}

func actionAddAsync(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Entrez le nouveau mot : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	_, err := d.Get(word)
	if err == nil {
		fmt.Printf("Le mot '%s' existe déjà dans le dictionnaire.\n", word)
		return
	}

	fmt.Print("Entrez la nouvelle définition : ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.AddAsync(word, definition)

	fmt.Printf("Le mot '%s' avec la définition '%s' a été ajouté.\n", word, definition)
}

func actionDefineAsync(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Entrez le mot : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	_, err := d.Get(word)
	if err != nil {
		fmt.Printf("Le mot '%s' n'existe pas dans le dico.\n", word)
		return
	}

	fmt.Print("Entrez la nouvelle définition : ")
	newDefinition, _ := reader.ReadString('\n')
	newDefinition = strings.TrimSpace(newDefinition)

	d.RemoveAsync(word)
	d.AddAsync(word, newDefinition)

	fmt.Printf("La définition pour le mot '%s' a été mise à jour.\n", word)
}

func actionRemoveAsync(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Écrivez le mot à supprimer : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.RemoveAsync(word)
}

func actionList(d *dictionary.Dictionary) {
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
