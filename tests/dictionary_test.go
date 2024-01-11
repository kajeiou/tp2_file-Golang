package tests

import (
	"os"
	"testing"

	"tp2/dictionary"

	"github.com/stretchr/testify/assert"
)

func TestDictionary_AddAsync(t *testing.T) {
	d := dictionary.New("test_data.csv")
	t.Log("Création du fichier csv de test")

	d.AddAsync("test1", "definition1")
	t.Log("Test ajout de mot")

	word, err := d.Get("test1")
	assert.NoError(t, err)
	assert.Equal(t, "test1", word.Word)
	assert.Equal(t, "definition1", word.Definition)

	d.AddAsync("test1", "definition1")

	words := d.List()

	assert.Equal(t, 1, len(words), "Le nombre de mots ne correspond pas")
	os.Remove("test_data.csv")
}

func TestDictionary_EditAsync(t *testing.T) {
	d := dictionary.New("test_data.csv")
	t.Log("Création du fichier csv de test")

	d.AddAsync("test1", "definition1")
	t.Log("Mot ajouté avec succès")

	d.EditAsync("test1", "nouvelle_definition")
	t.Log("Test de modification de mot")

	word, err := d.Get("test1")
	assert.NoError(t, err)
	assert.Equal(t, "nouvelle_definition", word.Definition)

	os.Remove("test_data.csv")
}

func TestDictionary_RemoveAsync(t *testing.T) {
	d := dictionary.New("test_data.csv")
	t.Log("Création du fichier csv de test")

	d.AddAsync("test1", "definition1")
	t.Log("Mot ajouté avec succès")

	d.RemoveAsync("test1")
	t.Log("Test de suppression de mot")

	words := d.List()
	assert.Equal(t, 0, len(words), "Le dictionnaire devrait être vide après la suppression")

	os.Remove("test_data.csv")
}
