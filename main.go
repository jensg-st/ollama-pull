package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/jensg-st/ollama-pull/server"
	"golang.org/x/crypto/ssh"
)

const (
	regirtsy = "registry.ollama.ai"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("no model name provided")
	}

	model := os.Args[1]

	slog.Info(fmt.Sprintf("using ollama model %s", model))

	dir, err := createOllamaDir()
	if err != nil {
		log.Fatalf("can not create ollama home dir: %s", err.Error())
	}

	err = createKeyPair(dir)
	if err != nil {
		log.Fatalf("can not create ollama keys: %s", err.Error())
	}

	slog.Info(fmt.Sprintf("using ollama dir %s", dir))

	err = server.PullModel(context.Background(), model, &server.RegistryOptions{}, logProgress)
	if err != nil {
		log.Fatalf("can not download ollama model %s: %s", model, err.Error())
	}

	server.FixBlobs(path.Join(dir, "/models/blobs"))

}

func createOllamaDir() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".ollama"), nil

}

func createKeyPair(dir string) error {

	privKeyPath := filepath.Join(dir, "id_ed25519")
	pubKeyPath := filepath.Join(dir, "id_ed25519.pub")

	_, err := os.Stat(privKeyPath)
	if os.IsNotExist(err) {
		cryptoPublicKey, cryptoPrivateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}

		privateKeyBytes, err := ssh.MarshalPrivateKey(cryptoPrivateKey, "")
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(privKeyPath), 0o755); err != nil {
			return fmt.Errorf("could not create directory %w", err)
		}

		if err := os.WriteFile(privKeyPath, pem.EncodeToMemory(privateKeyBytes), 0o600); err != nil {
			return err
		}

		sshPublicKey, err := ssh.NewPublicKey(cryptoPublicKey)
		if err != nil {
			return err
		}

		publicKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)

		if err := os.WriteFile(pubKeyPath, publicKeyBytes, 0o644); err != nil {
			return err
		}
	}

	return nil

}

func logProgress(s server.ProgressResponse) {

	percentage := "0%"
	if s.Total > 0 {
		total := (float64(s.Completed) / float64(s.Total)) * 100
		percentage = fmt.Sprintf("%.2f%%", total)
	}

	slog.Info(fmt.Sprintf("%s %s/%s (%s)\n", s.Status, fileSize(float64(s.Completed)), fileSize(float64(s.Total)), percentage))

}

func fileSize(s float64) string {

	sizes := []string{"B", "kB", "MB", "GB", "TB"}

	unitsLimit := len(sizes)

	i := 0
	for s >= 1024.0 && i < unitsLimit {
		s = s / 1024.0
		i++
	}

	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}

	return fmt.Sprintf(f, s, sizes[i])
}
