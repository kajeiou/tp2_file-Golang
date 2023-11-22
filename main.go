package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"tp2/dictionary"
)

func main() {
	mode := getModeFromArgs()

	myDictionary := dictionary.New("dictionary.csv")
	fmt.Println("Bienvenue dans le dico !")

	switch mode {
	case "console", "1":
		runConsoleMode(myDictionary)
	case "api", "2":
		runAPIMode(myDictionary)
	default:
		fmt.Println("Mode non reconnu. Choisissez le mode :")
		fmt.Println("1. Console")
		fmt.Println("2. API")

		reader := bufio.NewReader(os.Stdin)
		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erreur de lecture de l'entrée utilisateur:", err)
			return
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			runConsoleMode(myDictionary)
		case "2":
			runAPIMode(myDictionary)
		default:
			fmt.Println("Choix invalide. Terminé.")
		}
	}
}

func getModeFromArgs() string {
	if len(os.Args) < 2 {
		return ""
	}
	return strings.ToLower(os.Args[1])
}

func runConsoleMode(d *dictionary.Dictionary) {
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
			actionList(d)
		case "2":
			actionAddAsync(d, reader)
		case "3":
			actionDefineAsync(d, reader)
		case "4":
			actionRemoveAsync(d, reader)
		case "5":
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Choix invalide. Veuillez entrer un numéro valide.")
		}
	}
}

func runAPIMode(d *dictionary.Dictionary) {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/api/words/add", apiAddWordHandler(d))
	http.HandleFunc("/api/words/define/", apiDefineWordHandler(d))
	http.HandleFunc("/api/words/remove/", apiRemoveWordHandler(d))
	http.HandleFunc("/api/words/list", apiListWordsHandler(d))

	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
