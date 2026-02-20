package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"path/filepath"
	"strings"
	"testing"
)

func generateTestKeyPair(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	return key
}

// hybridEncrypt performs RSA-OAEP + AES-GCM encryption (the client-side counterpart to DecryptPayload).
func hybridEncrypt(t *testing.T, plaintext []byte, pub *rsa.PublicKey) (ek, nonce, ct []byte) {
	t.Helper()
	// Generate random AES-256 key
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}
	// RSA-OAEP encrypt the AES key
	var err error
	ek, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, aesKey, nil)
	if err != nil {
		t.Fatalf("EncryptOAEP: %v", err)
	}
	// AES-GCM encrypt the plaintext
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		t.Fatalf("NewCipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("NewGCM: %v", err)
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		t.Fatalf("rand.Read nonce: %v", err)
	}
	ct = gcm.Seal(nil, nonce, plaintext, nil)
	return
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := generateTestKeyPair(t)
	plaintext := []byte(`{"tracker_id":"abc","visitor_id":"xyz"}`)

	ek, nonce, ct := hybridEncrypt(t, plaintext, &key.PublicKey)

	got, err := DecryptPayload(ek, nonce, ct, key)
	if err != nil {
		t.Fatalf("DecryptPayload: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Errorf("decrypted = %q, want %q", got, plaintext)
	}
}

func TestDecryptPayload_BadKey(t *testing.T) {
	key1 := generateTestKeyPair(t)
	key2 := generateTestKeyPair(t)
	plaintext := []byte("secret data")

	ek, nonce, ct := hybridEncrypt(t, plaintext, &key1.PublicKey)

	_, err := DecryptPayload(ek, nonce, ct, key2) // wrong key
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptPayload_CorruptedCiphertext(t *testing.T) {
	key := generateTestKeyPair(t)
	plaintext := []byte("test data")

	ek, nonce, ct := hybridEncrypt(t, plaintext, &key.PublicKey)

	// Corrupt the ciphertext
	ct[0] ^= 0xFF

	_, err := DecryptPayload(ek, nonce, ct, key)
	if err == nil {
		t.Fatal("expected error for corrupted ciphertext")
	}
}

func TestPublicKeyPEM(t *testing.T) {
	key := generateTestKeyPair(t)
	pemStr, err := PublicKeyPEM(&key.PublicKey)
	if err != nil {
		t.Fatalf("PublicKeyPEM: %v", err)
	}
	if !strings.HasPrefix(pemStr, "-----BEGIN PUBLIC KEY-----") {
		t.Errorf("PEM should start with BEGIN PUBLIC KEY header, got %q", pemStr[:40])
	}
	if !strings.Contains(pemStr, "-----END PUBLIC KEY-----") {
		t.Error("PEM should contain END PUBLIC KEY footer")
	}
}

func TestEnsureKeyPair(t *testing.T) {
	dir := t.TempDir()
	privPath := filepath.Join(dir, "keys", "priv.pem")
	pubPath := filepath.Join(dir, "keys", "pub.pem")

	// First call should create keys
	if err := EnsureKeyPair(privPath, pubPath); err != nil {
		t.Fatalf("EnsureKeyPair (first): %v", err)
	}

	// Second call should be a no-op (idempotent)
	if err := EnsureKeyPair(privPath, pubPath); err != nil {
		t.Fatalf("EnsureKeyPair (second): %v", err)
	}

	// Verify keys can be loaded
	_, _, err := LoadKeyPair(privPath, pubPath)
	if err != nil {
		t.Fatalf("LoadKeyPair after EnsureKeyPair: %v", err)
	}
}

func TestLoadKeyPair(t *testing.T) {
	dir := t.TempDir()
	privPath := filepath.Join(dir, "priv.pem")
	pubPath := filepath.Join(dir, "pub.pem")

	if err := EnsureKeyPair(privPath, pubPath); err != nil {
		t.Fatalf("EnsureKeyPair: %v", err)
	}

	priv, pub, err := LoadKeyPair(privPath, pubPath)
	if err != nil {
		t.Fatalf("LoadKeyPair: %v", err)
	}
	if priv == nil {
		t.Fatal("private key is nil")
	}
	if pub == nil {
		t.Fatal("public key is nil")
	}

	// Verify the loaded public key matches the private key's public key
	if priv.PublicKey.N.Cmp(pub.N) != 0 {
		t.Error("loaded public key does not match private key")
	}
}
