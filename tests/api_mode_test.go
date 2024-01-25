package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"tp2/api_mode"
	"tp2/db"
	"tp2/dictionary"

	"github.com/stretchr/testify/assert"
)

func TestWelcomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api_mode.WelcomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Bienvenue dans le dico !")
}

func loginAndGetToken(t *testing.T) string {
	loginRequest := map[string]string{"username": "nabil", "password": "10"}
	loginJSON, err := json.Marshal(loginRequest)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(loginJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(api_mode.LoginHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	token := rr.Body.String()
	fmt.Println("Token:", token)
	assert.NotEmpty(t, token, "Token not found in the response")

	isValid := api_mode.IsValidToken(token)
	assert.True(t, isValid, "Le jeton généré doit être valide.")

	return token
}

func TestAddHandler(t *testing.T) {
	token := loginAndGetToken(t)

	word := dictionary.Word{Word: "test", Definition: "definition"}
	wordJSON, err := json.Marshal(word)
	assert.NoError(t, err)

	addWordReq, err := http.NewRequest("POST", "/api/words/add", bytes.NewBuffer(wordJSON))
	assert.NoError(t, err)
	addWordReq.Header.Set("Authorization", token)

	addWordRR := httptest.NewRecorder()
	wordRepository := &db.GormWordRepository{}
	myDictionary := dictionary.New("dictionary.csv", wordRepository)
	addWordHandler := http.HandlerFunc(api_mode.ApiAddWordHandler(myDictionary))
	addWordHandler.ServeHTTP(addWordRR, addWordReq)

	assert.Equal(t, http.StatusCreated, addWordRR.Code)
	assert.Contains(t, addWordRR.Body.String(), fmt.Sprintf("Le mot '%s' avec la définition '%s' a été ajouté.", word.Word, word.Definition))
}
