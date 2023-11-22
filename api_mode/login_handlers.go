package api_mode

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	//"github.com/dgrijalva/jwt-go"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Printf("Bad request method: %s, expected POST. Route: %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "La méthode "+r.Method+" n'est pas autorisée pour cette route. Utilisez POST.")
		return
	}
}

func generateToken() (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if secretKey == nil {
		log.Fatal("La variable d'environnement SECRET_KEY n'est pas définie.")
	}
}
