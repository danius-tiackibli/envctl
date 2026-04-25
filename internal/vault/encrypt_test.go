package vault_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envctl/internal/vault"
)

func TestDeriveKey(t *testing.T) {
	key := vault.DeriveKey("supersecret")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	key2 := vault.DeriveKey("supersecret")
	if !bytes.Equal(key, key2) {
		t.Fatal("expected same passphrase to produce same key")
	}

	keyDiff := vault.DeriveKey("different")
	if bytes.Equal(key, keyDiff) {
		t.Fatal("expected different passphrases to produce different keys")
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := vault.DeriveKey("test-passphrase")
	plaintext := []byte("MY_SECRET=hello\nANOTHER_VAR=world")

	ciphertext, err := vault.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := vault.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key := vault.DeriveKey("correct-passphrase")
	wrongKey := vault.DeriveKey("wrong-passphrase")
	plaintext := []byte("SECRET=value")

	ciphertext, err := vault.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = vault.Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptTooShort(t *testing.T) {
	key := vault.DeriveKey("passphrase")
	_, err := vault.Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for too-short ciphertext")
	}
}
