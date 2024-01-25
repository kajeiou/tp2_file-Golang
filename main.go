package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"tp2/api_mode"
	"tp2/console_mode"
	"tp2/db"
	"tp2/dictionary"
)

func main() {

	// Créez une instance du repository GormWordRepository
	wordRepository := &db.GormWordRepository{}

	// Initialisez la base de données
	err := wordRepository.InitializeDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer wordRepository.CloseDB()

	// Ajoutez un mot à la base de données
	/*err = wordRepository.AddWordToDB("example", "This is an example definition.")
	if err != nil {
		log.Fatal("Failed to add word to database:", err)
	}*/

	mode := getModeFromArgs()

	myDictionary := dictionary.New("dictionary.csv", wordRepository)
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
			console_mode.ActionList(d)
		case "2":
			console_mode.ActionAddAsync(d, reader)
		case "3":
			console_mode.ActionDefineAsync(d, reader)
		case "4":
			console_mode.ActionRemoveAsync(d, reader)
		case "5":
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Choix invalide. Veuillez entrer un numéro valide.")
		}
	}
}

func runAPIMode(d *dictionary.Dictionary) {
	http.HandleFunc("/", api_mode.WelcomeHandler)
	http.HandleFunc("/api/words/add", api_mode.ApiAddWordHandler(d))
	http.HandleFunc("/api/words/define/", api_mode.ApiDefineWordHandler(d))
	http.HandleFunc("/api/words/remove/", api_mode.ApiRemoveWordHandler(d))
	http.HandleFunc("/api/words/list", api_mode.ApiListWordsHandler(d))
	http.HandleFunc("/api/login", api_mode.LoginHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Println("Starting server on", port)
	api_mode.LogToFile("runAPIMode", fmt.Sprintf("Server started on %s", port))
	log.Fatal(http.ListenAndServe(port, nil))
}
