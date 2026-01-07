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

	"github.com/hrz8/altalune/internal/config"
	"github.com/hrz8/altalune/internal/shared/crypto"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	secretFlag := flag.String("secret", "", "Raw client secret (optional)")
	flag.Parse()

	secret, err := resolveSecret(*secretFlag)
	if err != nil {
		log.Fatalf("failed to read client secret: %v", err)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	key := cfg.GetIAMEncryptionKey()
	if err := crypto.ValidateKey(key); err != nil {
		log.Fatalf("invalid encryption key: %v", err)
	}

	encrypted, err := crypto.Encrypt(secret, key)
	if err != nil {
		log.Fatalf("failed to encrypt client secret: %v", err)
	}

	fmt.Println(encrypted)
}

func resolveSecret(fromFlag string) (string, error) {
	if fromFlag != "" {
		return fromFlag, nil
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("inspect stdin: %w", err)
	}

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

	fmt.Fprint(os.Stderr, "Enter client secret: ")
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
