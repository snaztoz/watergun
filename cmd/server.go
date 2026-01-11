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
	serverFlagAdminKey      string

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server process",
		Long: `Start the main chat-server process.

It requires the "--public-key" option to be provided in order to
authenticate all incoming user connection requests. Currently, it
only supports the public key of ECDSA-256 (for JWT ES256).

On the other hand, the admin management secret key ("--admin-key"
flag) is not required. In this case, the server will use a static
key of value "ADMIN-INSECURE-KEY" to authenticate incoming admin
requests. Therefore, it is not suitable to be used in production
setting and a warning log message will be printed to the screen
every time the server runs.
`,
		Run: func(_ *cobra.Command, _ []string) {
			runServer()
		},
	}
)

func init() {
	serverCmd.Flags().StringVarP(&serverFlagPublicKeyPath, "public-key", "P", "", "public-key file path (required)")
	serverCmd.Flags().StringVarP(&serverFlagAdminKey, "admin-key", "a", "", "admin key. Default: \"ADMIN-INSECURE-KEY\"")

	serverCmd.MarkFlagRequired("public-key")

	rootCmd.AddCommand(serverCmd)
}

func runServer() {
	log.SetFlags(0)

	if serverFlagAdminKey == "" {
		log.Printf("WARNING: the current admin key is *not secure* and *should not* be used in production\n\n")
	}

	s := server.New(port, serverFlagAdminKey, readPublicKeyPEMFile())

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
