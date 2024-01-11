package api_mode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		LogToFile("LoginHandler", fmt.Sprintf("Mauvaise méthode de requête: %s, POST attendu. Route: %s", r.Method, r.URL.Path))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez POST.")
		return
	}

	var requestBody map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Erreur lors de la lecture du corps de la requête.")
		return
	}

	username, usernameExists := requestBody["username"]
	password, passwordExists := requestBody["password"]

	if !usernameExists || !passwordExists {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Nom d'utilisateur et mot de passe requis.")
		return
	}

	if !isValidUser(username, password) {
		LogToFile("LoginHandler", fmt.Sprintf("Nom d'utilisateur ou mot de passe incorrect pour l'utilisateur: %s", username))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Nom d'utilisateur ou mot de passe incorrect.")
		return
	}

	token, err := generateToken(username)
	if err != nil {
		LogToFile("LoginHandler", fmt.Sprintf("Erreur lors de la génération du jeton d'authentification pour l'utilisateur: %s", username))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Erreur lors de la génération du jeton d'authentification.")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, token)
}

func isValidUser(username, password string) bool {
	return username == "nabil" && password == "10"
}

func generateToken(username string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if secretKey == nil {
		log.Fatal("La variable d'environnement SECRET_KEY n'est pas définie.")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
