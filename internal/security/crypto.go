package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func LoadKeyPair(privPath, pubPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privPEM, err := os.ReadFile(privPath)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, nil, errors.New("failed to decode private key PEM")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	pubPEM, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, nil, err
	}
	pubBlock, _ := pem.Decode(pubPEM)
	if pubBlock == nil {
		return nil, nil, errors.New("failed to decode public key PEM")
	}
	pubIface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	pubKey, ok := pubIface.(*rsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("not an RSA public key")
	}

	return privKey, pubKey, nil
}

// DecryptPayload decrypts a hybrid RSA+AES-GCM encrypted payload.
// ek is the RSA-OAEP encrypted AES key, nonce is the AES-GCM nonce, ct is the ciphertext.
func DecryptPayload(ek, nonce, ct []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	dataKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ek, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt data key")
	}

	block, err := aes.NewCipher(dataKey)
	if err != nil {
		return nil, errors.New("failed to create AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM")
	}

	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt payload")
	}

	return plaintext, nil
}

func PublicKeyPEM(pub *rsa.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return string(pem.EncodeToMemory(block)), nil
}
