package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

func EnsureKeyPair(privPath, pubPath string) error {
	_, errPriv := os.Stat(privPath)
	_, errPub := os.Stat(pubPath)
	if errPriv == nil && errPub == nil {
		return nil // both files exist
	}

	if err := os.MkdirAll(filepath.Dir(privPath), 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(pubPath), 0700); err != nil {
		return err
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	if err := os.WriteFile(privPath, privPEM, 0600); err != nil {
		return err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	if err := os.WriteFile(pubPath, pubPEM, 0644); err != nil {
		return err
	}

	return nil
}
