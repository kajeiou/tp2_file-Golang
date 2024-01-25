package api_mode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tp2/dictionary"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	LogAndRespond(w, r, "Bienvenue dans le dico !", http.StatusOK)
}

func ApiAddWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authenticateRequest(w, r) {
			return
		}

		if r.Method != http.MethodPost {
			logMessage := fmt.Sprintf("Mauvaise méthode de requête : %s, attendue POST %s", r.Method, r.URL.Path)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		var word dictionary.Word
		err := json.NewDecoder(r.Body).Decode(&word)
		if err != nil {
			logMessage := fmt.Sprintf("Error decoding request body: %v. Route: %s", err, r.URL.Path)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		if word.Word == "" || word.Definition == "" {
			logMessage := fmt.Sprintf("Clés manquantes dans le corps de la requête. Route: %s", r.URL.Path)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		if err := validateWordAndDefinitionLength(word.Word, word.Definition); err != nil {
			logMessage := fmt.Sprintf("Erreur de validation : %v", err)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		if err := d.AddAsync(word.Word, word.Definition); err != nil {
			logMessage := fmt.Sprintf("Erreur lors de l'ajout du mot : %v", err)
			LogAndRespond(w, r, logMessage, http.StatusInternalServerError)
			return
		}

		LogAndRespond(w, r, fmt.Sprintf("Le mot '%s' avec la définition '%s' a été ajouté.", word.Word, word.Definition), http.StatusCreated)
	}
}

func ApiDefineWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authenticateRequest(w, r) {
			return
		}
		if r.Method != http.MethodPut {
			logMessage := fmt.Sprintf("Mauvaise méthode de requête : %s, PUT attendu. Route: %s", r.Method, r.URL.Path)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		word := extractWordFromURL(r.URL.Path)

		if word == "" {
			LogAndRespond(w, r, "Veuillez saisir un mot dans l'URL.", http.StatusBadRequest)
			return
		}

		var newDefinition string
		err := json.NewDecoder(r.Body).Decode(&newDefinition)
		if err != nil {
			LogAndRespond(w, r, "Corps de la demande non valide", http.StatusBadRequest)
			return
		}

		if err := validateWordAndDefinitionLength(word, newDefinition); err != nil {
			logMessage := fmt.Sprintf("Erreur de validation : %v", err)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		err = d.EditAsync(word, newDefinition)
		if err != nil {
			logMessage := fmt.Sprintf("Erreur lors de la mise à jour de la définition dans la base de données : %v", err)
			LogAndRespond(w, r, logMessage, http.StatusInternalServerError)
			return
		}

		LogAndRespond(w, r, fmt.Sprintf("La définition pour le mot '%s' a été mise à jour.", word), http.StatusOK)
	}
}

func ApiRemoveWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authenticateRequest(w, r) {
			return
		}

		if r.Method != http.MethodDelete {
			logMessage := fmt.Sprintf("Mauvaise méthode de requête: %s, Delete attendu. Route: %s", r.Method, r.URL.Path)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		word := extractWordFromURL(r.URL.Path)

		if word == "" {
			LogAndRespond(w, r, "Veuillez saisir un mot dans l'URL.", http.StatusBadRequest)
			return
		}

		err := d.RemoveAsync(word)
		if err != nil {
			logMessage := fmt.Sprintf("Erreur lors de la suppression du mot dans la base de données : %v", err)
			LogAndRespond(w, r, logMessage, http.StatusInternalServerError)
			return
		}

		logMessage := fmt.Sprintf("Suppression du mot %s", word)
		LogAndRespond(w, r, logMessage, http.StatusOK)
	}
}

func ApiListWordsHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authenticateRequest(w, r) {
			return
		}

		if r.Method != http.MethodGet {
			logMessage := fmt.Sprintf("Mauvaise méthode de requête :%s, GET attendu. Route: %s", r.Method, r.URL.Path)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		wordsList, err := d.List()
		if err != nil {
			LogAndRespond(w, r, fmt.Sprintf("Erreur lors de la récupération de la liste des mots : %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if len(wordsList) == 0 {
			LogAndRespond(w, r, "Aucun mot dans le dico.", http.StatusOK)
		} else {
			logMessage := fmt.Sprintf("Requête : %s. Route: %s", r.Method, r.URL.Path)
			LogToFile("ApiListWordsHandler", logMessage)
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
