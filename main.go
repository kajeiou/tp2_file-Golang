package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tp1/dictionary"
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
			actionAdd(myDictionary, reader)
		case "3":
			actionDefine(myDictionary, reader)
		case "4":
			actionRemove(myDictionary, reader)
		case "5":
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Choix invalide. Veuillez entrer un numéro valide.")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Entrez le nouveau mot : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	_, exists := d.Get(word)
	if exists {
		fmt.Printf("Le mot '%s' existe déjà dans le dictionnaire.\n", word)
		return
	}

	fmt.Print("Entrez la nouvelle définition : ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)

	fmt.Printf("Le mot '%s' avec la définition '%s' a été ajouté.\n", word, definition)
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Entrez le mot : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Entrez la nouvelle définition : ")
	newDefinition, _ := reader.ReadString('\n')
	newDefinition = strings.TrimSpace(newDefinition)

	d.Remove(word)

	d.Add(word, newDefinition)

	fmt.Printf("La définition pour le mot '%s' a été mise à jour.\n", word)
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {

	fmt.Print("Écrivez le mot à supprimer : ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)
}

func actionList(d *dictionary.Dictionary) {
	wordsList, _ := d.List()
	if len(wordsList) == 0 {
		fmt.Println("Aucun mot dans le dico.")
	} else {
		fmt.Println("Liste des mots du dico:")
		for _, word := range wordsList {
			fmt.Printf("%s\n", word)
		}
	}
}
