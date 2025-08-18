package nanoid

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

func GeneratePublicIDBatch(number int) []string {
	results := make([]string, 0)
	for i := 0; i < number; i++ {
		id, _ := GeneratePublicID()
		results = append(results, id)
	}

	return results
}
