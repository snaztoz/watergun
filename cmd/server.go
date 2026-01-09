package cmd

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/snaztoz/watergun/server"
)

const port = "8080"

var (
	serverFlagPublicKeyPath string

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server process",
		Run: func(_ *cobra.Command, _ []string) {
			runServer()
		},
	}
)

func init() {
	serverCmd.Flags().StringVarP(&serverFlagPublicKeyPath, "public-key", "P", "", "public-key file path (required)")
	serverCmd.MarkFlagRequired("public-key")

	rootCmd.AddCommand(serverCmd)
}

func runServer() {
	log.SetFlags(0)

	s := server.New(port, readPublicKeyPEMFile())

	go s.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.Stop()
}

func readPublicKeyPEMFile() crypto.PublicKey {
	pemContent, err := os.ReadFile(serverFlagPublicKeyPath)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	block, _ := pem.Decode(pemContent)
	if block == nil {
		log.Fatal("failed to parse the PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse public key: %v", err)
	}

	key, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("unsupported public key type")
	}

	return key
}
