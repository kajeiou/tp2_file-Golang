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
	LogToFile("WelcomeHandler", fmt.Sprintf("Requête reçue : %s %s", r.Method, r.URL.Path))
	fmt.Fprintln(w, "Bienvenue dans le dico !")
}

func authenticateRequest(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		LogToFile("DictionnaryHandler", fmt.Sprintf("Requête non autorisé : %s %s. Le jeton d'authentification est requis.", r.Method, r.URL.Path))
		fmt.Fprint(w, "Accès non autorisé. Le jeton d'authentification est requis.")
		return false
	}

	if !isValidToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		LogToFile("DictionnaryHandler", fmt.Sprintf("Requête non autorisé : %s %s. Le jeton d'authentification est invalide.", r.Method, r.URL.Path))
		fmt.Fprint(w, "Accès non autorisé. Jeton d'authentification invalide.")
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
		//log.Printf("Erreur lors de la validation du jeton: %v", err)
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
			LogToFile("ApiAddWordHandler", logMessage)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez POST.")
			return
		}

		var word dictionary.Word
		err := json.NewDecoder(r.Body).Decode(&word)
		if err != nil {
			logMessage := fmt.Sprintf("Error decoding request body: %v. Route: %s", err, r.URL.Path)
			LogToFile("ApiAddWordHandler", logMessage)
			http.Error(w, "Corps de la demande non valide", http.StatusBadRequest)
			return
		}

		if word.Word == "" || word.Definition == "" {
			logMessage := fmt.Sprintf("Clés manquantes dans le corps de la requête. Route: %s", r.URL.Path)
			LogToFile("ApiAddWordHandler", logMessage)
			http.Error(w, "Le corps de la requête doit contenir les clés 'Word' et 'Definition'.", http.StatusBadRequest)
			return
		}

		_, err = d.Get(word.Word)
		if err == nil {
			logMessage := fmt.Sprintf("Conflict: Le mot '%s' existe déjà dans le dico. Route: %s", word.Word, r.URL.Path)
			LogToFile("ApiAddWordHandler", logMessage)
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

		if !authenticateRequest(w, r) {
			return
		}
		if r.Method != http.MethodPut {
			log.Printf("Mauvaise méthode de requête : %s, PUT attendu. Route: %s", r.Method, r.URL.Path)
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
			log.Printf("Erreur de décodage du corps de la requête : %v. Route: %s", err, r.URL.Path)
			http.Error(w, "Corps de la demande non valide", http.StatusBadRequest)
			return
		}

		d.EditAsync(existingWord.Word, newDefinition)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "La définition pour le mot '%s' a été mise à jour.", word)
	}
}

func ApiRemoveWordHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !authenticateRequest(w, r) {
			return
		}

		if r.Method != http.MethodDelete {
			log.Printf("Mauvaise méthode de requête: %s, Delete attendu. Route: %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez Delete.")
			return
		}

		word := extractWordFromURL(r.URL.Path)

		if word == "" {
			log.Print("Mauvaise demande: Veuillez saisir un mot dans l'URL.")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Veuillez saisir un mot dans l'URL.")
			return
		}

		log.Printf("Suppression du mot %s", word)
		d.RemoveAsync(word)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Le mot '%s' a été supprimé.", word)
	}
}

func ApiListWordsHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !authenticateRequest(w, r) {
			return
		}

		if r.Method != http.MethodGet {
			log.Printf("Mauvaise méthode de requête :%s, GET attendu. Route: %s", r.Method, r.URL.Path)
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
