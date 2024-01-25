package api_mode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if !isTesting() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Erreur lors du chargement du fichier .env")
		}
	}
}

func isTesting() bool {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			return true
		}
	}
	return false
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		LogAndRespond(w, r, fmt.Sprintf("Mauvaise méthode de requête: %s, POST attendu. Route: %s", r.Method, r.URL.Path), http.StatusBadRequest)
		return
	}

	var requestBody map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		LogAndRespond(w, r, "Erreur lors de la lecture du corps de la requête.", http.StatusBadRequest)
		return
	}

	username, usernameExists := requestBody["username"]
	password, passwordExists := requestBody["password"]

	if !usernameExists || !passwordExists {
		LogAndRespond(w, r, "Nom d'utilisateur et mot de passe requis.", http.StatusBadRequest)
		return
	}

	if !isValidUser(username, password) {
		LogAndRespond(w, r, fmt.Sprintf("Nom d'utilisateur ou mot de passe incorrect pour l'utilisateur: %s", username), http.StatusUnauthorized)
		return
	}

	token, err := generateToken(username)
	if err != nil {
		LogAndRespond(w, r, fmt.Sprintf("Erreur lors de la génération du jeton d'authentification pour l'utilisateur: %s", username), http.StatusInternalServerError)
		return
	}

	LogAndRespond(w, r, token, http.StatusOK)
}

func isValidUser(username, password string) bool {
	return username == "nabil" && password == "10"
}
