package api_mode

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

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
