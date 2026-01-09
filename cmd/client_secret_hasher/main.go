package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/hrz8/altalune/internal/shared/password"
)

// Production parameters for Argon2id hashing
const (
	defaultIterations = 2         // Time cost (t)
	defaultMemory     = 64 * 1024 // Memory cost in KB (m) = 64MB
	defaultThreads    = 4         // Parallelism (p)
	defaultLength     = 32        // Hash length in bytes
)

func main() {
	// CLI flags
	secretFlag := flag.String("secret", "", "Client secret to hash (optional, can use stdin)")
	iterations := flag.Uint("iterations", defaultIterations, "Time cost (iterations)")
	memory := flag.Uint("memory", defaultMemory, "Memory cost in KB")
	threads := flag.Uint("threads", defaultThreads, "Parallelism (threads)")
	length := flag.Uint("length", defaultLength, "Hash length in bytes")
	flag.Parse()

	// Resolve secret from flag or stdin
	secret, err := resolveSecret(*secretFlag)
	if err != nil {
		log.Fatalf("failed to read client secret: %v", err)
	}

	// Validate minimum secret length (32 characters for OAuth client secrets)
	if len(secret) < 32 {
		log.Fatalf("client secret must be at least 32 characters, got %d", len(secret))
	}

	// Hash the secret with Argon2id
	hashedSecret, err := password.HashPassword(secret, password.HashOption{
		Iterations: uint32(*iterations),
		Memory:     uint32(*memory),
		Threads:    uint8(*threads),
		Len:        uint32(*length),
	})
	if err != nil {
		log.Fatalf("failed to hash client secret: %v", err)
	}

	// Output the PHC string format hash
	fmt.Println(hashedSecret)
}

// resolveSecret resolves the secret from flag, stdin pipe, or interactive prompt
func resolveSecret(fromFlag string) (string, error) {
	// If provided via flag, use it
	if fromFlag != "" {
		return fromFlag, nil
	}

	// Check if stdin is a pipe (non-interactive)
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("inspect stdin: %w", err)
	}

	// If stdin is a pipe, read from it
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(bufio.NewReader(os.Stdin))
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		secret := strings.TrimRight(string(data), "\r\n")
		if secret == "" {
			return "", errors.New("stdin provided empty client secret")
		}
		return secret, nil
	}

	// Interactive mode: prompt for secret
	fmt.Fprint(os.Stderr, "Enter client secret (min 32 chars): ")
	reader := bufio.NewReader(os.Stdin)
	secret, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read input: %w", err)
	}
	secret = strings.TrimRight(secret, "\r\n")
	if secret == "" {
		return "", errors.New("client secret cannot be empty")
	}
	return secret, nil
}
