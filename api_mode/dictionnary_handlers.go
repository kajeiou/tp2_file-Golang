package api_mode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"tp2/dictionary"

	"github.com/dgrijalva/jwt-go"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	LogAndRespond(w, r, "Bienvenue dans le dico !", http.StatusOK)
}

func authenticateRequest(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		LogAndRespond(w, r, "Accès non autorisé. Le jeton d'authentification est requis.", http.StatusUnauthorized)
		return false
	}

	if !isValidToken(token) {
		LogAndRespond(w, r, "Accès non autorisé. Jeton d'authentification invalide.", http.StatusUnauthorized)
		return false
	}

	return true
}

func isValidToken(tokenString string) bool {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if secretKey == nil {
		log.Println("La variable d'environnement SECRET_KEY n'est pas définie.")
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Méthode de signature invalide: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
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

		_, err = d.Get(word.Word)
		if err == nil {
			logMessage := fmt.Sprintf("Conflict: Le mot '%s' existe déjà dans le dico. Route: %s", word.Word, r.URL.Path)
			LogAndRespond(w, r, logMessage, http.StatusConflict)
			return
		}

		if err := validateWordAndDefinitionLength(word.Word, word.Definition); err != nil {
			logMessage := fmt.Sprintf("Erreur de validation : %v", err)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		d.AddAsync(word.Word, word.Definition)
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

		existingWord, err := d.Get(word)
		if err != nil {
			LogAndRespond(w, r, fmt.Sprintf("Le mot '%s' n'existe pas dans le dico.", word), http.StatusNotFound)
			return
		}

		var newDefinition string
		err = json.NewDecoder(r.Body).Decode(&newDefinition)
		if err != nil {
			LogAndRespond(w, r, "Corps de la demande non valide", http.StatusBadRequest)
			return
		}

		if err := validateWordAndDefinitionLength(existingWord.Word, newDefinition); err != nil {
			logMessage := fmt.Sprintf("Erreur de validation : %v", err)
			LogAndRespond(w, r, logMessage, http.StatusBadRequest)
			return
		}

		d.EditAsync(existingWord.Word, newDefinition)
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

		if removed := d.RemoveAsync(word); !removed {
			logMessage := fmt.Sprintf("Le mot '%s' n'a pas été trouvé dans le dictionnaire", word)
			LogAndRespond(w, r, logMessage, http.StatusNotFound)
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

		wordsList := d.List()
		if len(wordsList) == 0 {
			LogAndRespond(w, r, "Aucun mot dans le dico.", http.StatusOK)
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
