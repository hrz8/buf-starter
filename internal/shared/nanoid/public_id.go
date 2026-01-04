package nanoid

import (
	"fmt"
)

const (
	PublicIDAlphabet = "23456789abcdefghjklmnpqrstuvwxyz"
	PublicIDSize     = 14
)

func GeneratePublicID() (string, error) {
	id, err := GenerateString(PublicIDAlphabet, PublicIDSize)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GeneratePublicIDBatch(number int) ([]string, error) {
	if number <= 0 {
		return nil, fmt.Errorf("number must be greater than 0")
	}

	results := make([]string, 0, number)

	for i := 0; i < number; i++ {
		id, err := GeneratePublicID()
		if err != nil {
			return nil, fmt.Errorf("failed generating ID at index %d: %w", i, err)
		}
		results = append(results, id)
	}

	return results, nil
}
