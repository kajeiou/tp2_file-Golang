package tests

import (
	"log"
	"testing"
	"tp2/db"
)

func TestCRUDOperations(t *testing.T) {
	// Initialisez la base de données en mémoire
	wordRepository := &db.GormWordRepository{}

	err := wordRepository.InitializeDB(":memory:")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer wordRepository.CloseDB()

	// Test d'ajout
	err = wordRepository.AddWordToDB("example", "This is an example definition.")
	if err != nil {
		log.Fatal("Failed to add word to database:", err)
	}

	// Test de récupération du mot ajouté
	word, err := wordRepository.GetWordFromDB("example")
	if err != nil {
		t.Errorf("Erreur lors de la récupération du mot ajouté : %v", err)
	}
	if word.Word != "example" || word.Definition != "This is an example definition." {
		t.Errorf("Le mot récupéré ne correspond pas à celui ajouté.")
	}

	// Test de modification
	err = wordRepository.UpdateWordInDB("example", "nouvelle_definition")
	if err != nil {
		t.Errorf("Erreur lors de la modification du mot : %v", err)
	}

	// Vérification de la modification
	word, err = wordRepository.GetWordFromDB("example")
	if err != nil {
		t.Errorf("Erreur lors de la récupération du mot modifié : %v", err)
	}
	if word.Definition != "nouvelle_definition" {
		t.Errorf("La définition récupérée ne correspond pas à celle modifiée.")
	}

	// Test de suppression
	err = wordRepository.DeleteWordFromDB("example")
	if err != nil {
		t.Errorf("Erreur lors de la suppression du mot : %v", err)
	}

	// Vérification de la suppression
	_, err = wordRepository.GetWordFromDB("example")
	if err == nil {
		t.Errorf("Le mot supprimé est toujours présent dans la base de données.")
	}
}
