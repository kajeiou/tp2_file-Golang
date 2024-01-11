package api_mode

import "fmt"

func validateWordAndDefinitionLength(word, definition string) error {
	minWordLength := 2
	maxWordLength := 30
	minDefinitionLength := 5
	maxDefinitionLength := 255

	if len(word) < minWordLength || len(word) > maxWordLength {
		return fmt.Errorf("La longueur du mot doit être entre %d et %d caractères", minWordLength, maxWordLength)
	}

	if len(definition) < minDefinitionLength || len(definition) > maxDefinitionLength {
		return fmt.Errorf("La longueur de la définition doit être entre %d et %d caractères", minDefinitionLength, maxDefinitionLength)
	}

	return nil
}
