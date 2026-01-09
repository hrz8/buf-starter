package password

import (
	"errors"
)

var (
	// ErrInvalidHashedString is returned when the encoded hash is not in the correct format
	ErrInvalidHashedString = errors.New("the encoded hash is not in the correct format")
	// ErrIncompatibleVersion is returned when the argon2 version is incompatible
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

const (
	// saltLen is the length of the salt in bytes
	saltLen = 16
	// threads is the number of threads to use for hashing (parallelism)
	threads = 4
	// splitN is the number of parts in the PHC string format
	splitN = 6
)

// HashOption contains the parameters for argon2id hashing
type HashOption struct {
	Iterations uint32 // Time cost (t)
	Memory     uint32 // Memory cost in KB (m)
	Threads    uint8  // Parallelism (p)
	Len        uint32 // Hash length in bytes
}
