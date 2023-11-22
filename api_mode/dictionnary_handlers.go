package api_mode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"tp2/dictionary"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	LogToFile(fmt.Sprintf("Received request: %s %s", r.Method, r.URL.Path))
	fmt.Fprintln(w, "Bienvenue dans le dico !")
}

func ApiAddWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Bad request method: %s, expected POST. Route: %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez POST.")
			return
		}

		var word dictionary.Word
		err := json.NewDecoder(r.Body).Decode(&word)
		if err != nil {
			log.Printf("Error decoding request body: %v. Route: %s", err, r.URL.Path)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err = d.Get(word.Word)
		if err == nil {
			log.Printf("Conflict: Le mot '%s' existe déjà dans le dico. Route: %s", word.Word, r.URL.Path)
			http.Error(w, fmt.Sprintf("Le mot '%s' existe déjà dans le dico.", word.Word), http.StatusConflict)
			return
		}

		d.AddAsync(word.Word, word.Definition)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Le mot '%s' avec la définition '%s' a été ajouté.", word.Word, word.Definition)
	}
}

func ApiDefineWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			log.Printf("Bad request method: %s, expected PUT. Route: %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez PUT.")
			return
		}

		word := extractWordFromURL(r.URL.Path)

		if word == "" {
			log.Print("Bad request: Veuillez saisir un mot dans l'URL.")
			http.Error(w, "Veuillez saisir un mot dans l'URL.", http.StatusBadRequest)
			return
		}

		existingWord, err := d.Get(word)
		if err != nil {
			log.Printf("Not Found: Le mot '%s' n'existe pas dans le dico. Route: %s", word, r.URL.Path)
			http.Error(w, fmt.Sprintf("Le mot '%s' n'existe pas dans le dico.", word), http.StatusNotFound)
			return
		}

		var newDefinition string
		err = json.NewDecoder(r.Body).Decode(&newDefinition)
		if err != nil {
			log.Printf("Error decoding request body: %v. Route: %s", err, r.URL.Path)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		d.EditAsync(existingWord.Word, newDefinition)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "La définition pour le mot '%s' a été mise à jour.", word)
	}
}

func ApiRemoveWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			log.Printf("Bad request method: %s, expected Delete. Route: %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez Delete.")
			return
		}

		word := extractWordFromURL(r.URL.Path)

		if word == "" {
			log.Print("Bad request: Veuillez saisir un mot dans l'URL.")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Veuillez saisir un mot dans l'URL.")
			return
		}

		log.Printf("Removing word: %s", word)
		d.RemoveAsync(word)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Le mot '%s' a été supprimé.", word)
	}
}

func ApiListWordsHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("Bad request method: %s, expected GET. Route: %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez GET.")
			return
		}

		wordsList := d.List()
		if len(wordsList) == 0 {
			log.Print("Aucun mot dans le dico.")
			fmt.Fprintln(w, "Aucun mot dans le dico.")
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(wordsList)
		}
	}
}

func extractWordFromURL(urlPath string) string {
	parts := strings.Split(urlPath, "/")
	if len(parts) == 5 {
		return parts[4]
	}
	return ""
}
