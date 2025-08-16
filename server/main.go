package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const (
	// DefaultHost is the default host to listen on
	DefaultHost = "0.0.0.0"
	// DefaultPort is the default port to listen on
	DefaultPort = "2222"
)

func main() {
	host := getEnvOrDefault("HOST", DefaultHost)
	port := getEnvOrDefault("PORT", DefaultPort)
	addr := fmt.Sprintf("%s:%s", host, port)

	// Generate a new private key for the server if it doesn't exist
	privateKey, err := getOrGeneratePrivateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKey(privateKey),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

	done := make(chan os.Signal, 1)
	log.Printf("Starting SSH server on %s", addr)
	if err = s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Get the path to the chat client binary
	clientPath, err := filepath.Abs("../tui/chat")
	if err != nil {
		log.Fatalf("Failed to get absolute path to chat client: %v", err)
	}

	// Check if the binary exists, if not, build it
	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		log.Printf("Building chat client...")
		cmd := exec.Command("go", "build", "-o", clientPath, "../tui/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to build chat client: %v", err)
		}
	}

	// Create a PTY for the chat client
	pty, _, _ := s.Pty()

	// Set up the command to run the chat client
	cmd := exec.Command(clientPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("TERM=%s", pty.Term))
	cmd.Stdin = s
	cmd.Stdout = s
	cmd.Stderr = s

	// Start the chat client
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start chat client: %v", err)
		fmt.Fprintf(s, "Failed to start chat client: %v\r\n", err)
		s.Exit(1)
	}

	// Wait for the chat client to exit
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Chat client exited with error: %v", err)
		}
		s.Close()
	}()

	// Return a dummy model since we're not using the Bubble Tea program directly
	// The actual chat client will handle the TUI
	return nil, nil
}

func getOrGeneratePrivateKey() (ed25519.PrivateKey, error) {
	keyPath := "ssh_host_ed25519"

	// Try to load existing key
	if _, err := os.Stat(keyPath); err == nil {
		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}

		privateKey, err := gossh.ParseRawPrivateKey(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		return privateKey.(ed25519.PrivateKey), nil
	}

	// Generate a new key
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Save the key for future use
	privateKeyPEM, err := gossh.MarshalPrivateKey(privKey, "")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := os.WriteFile(keyPath, privateKeyPEM, 0600); err != nil {
		return nil, fmt.Errorf("failed to write private key: %w", err)
	}

	log.Printf("Generated new SSH host key: %s", keyPath)
	return privKey, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
