package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tp2/api_mode"
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

func TestApiAddWordHandler(t *testing.T) {
	// Création du dictionnaire
	d := dictionary.New("test_data.csv")

	// Création d'une nouvelle requête POST avec un mot et une définition JSON dans le corps
	word := dictionary.Word{Word: "test", Definition: "definition"}
	wordJSON, err := json.Marshal(word)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/word", bytes.NewBuffer(wordJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api_mode.ApiAddWordHandler(d))

	handler.ServeHTTP(rr, req)

	// Vérifications
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "Le mot 'test' avec la définition 'definition' a été ajouté.")
}

func TestApiDefineWordHandler(t *testing.T) {
	// Création du dictionnaire
	d := dictionary.New("test_data.csv")

	// Ajout d'un mot au dictionnaire
	d.AddAsync("test", "definition")

	// Création d'une nouvelle requête PUT avec une nouvelle définition JSON dans le corps
	newDefinition := "nouvelle_definition"
	newDefinitionJSON, err := json.Marshal(newDefinition)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/word/test", bytes.NewBuffer(newDefinitionJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api_mode.ApiDefineWordHandler(d))

	handler.ServeHTTP(rr, req)

	// Vérifications
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "La définition pour le mot 'test' a été mise à jour.")
}

func TestApiRemoveWordHandler(t *testing.T) {
	// Création du dictionnaire
	d := dictionary.New("test_data.csv")

	// Ajout d'un mot au dictionnaire
	d.AddAsync("test", "definition")

	// Création d'une nouvelle requête DELETE
	req, err := http.NewRequest("DELETE", "/api/word/test", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api_mode.ApiRemoveWordHandler(d))

	handler.ServeHTTP(rr, req)

	// Vérifications
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Suppression du mot test")
}

func TestApiListWordsHandler(t *testing.T) {
	// Création du dictionnaire
	d := dictionary.New("test_data.csv")

	// Création d'une nouvelle requête GET
	req, err := http.NewRequest("GET", "/api/words", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api_mode.ApiListWordsHandler(d))

	handler.ServeHTTP(rr, req)

	// Vérifications
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Aucun mot dans le dico.")
}

func authenticateRequest(w http.ResponseWriter, r *http.Request) bool {
	// Fonction de simulation de l'authentification, toujours renvoie true
	return true
}
