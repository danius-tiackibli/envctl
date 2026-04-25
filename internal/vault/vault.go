package vault

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Vault represents an encrypted collection of environment variables.
type Vault struct {
	Name string
	Path string
}

// EnvMap is a map of environment variable key-value pairs.
type EnvMap map[string]string

// New creates a new Vault instance with the given name and storage directory.
func New(name, dir string) *Vault {
	return &Vault{
		Name: name,
		Path: filepath.Join(dir, name+".vault"),
	}
}

// Save encrypts the given EnvMap with the provided passphrase and writes it to disk.
func (v *Vault) Save(env EnvMap, passphrase string) error {
	if len(passphrase) == 0 {
		return errors.New("passphrase must not be empty")
	}

	data, err := json.Marshal(env)
	if err != nil {
		return err
	}

	key, salt, err := DeriveKey(passphrase, nil)
	if err != nil {
		return err
	}

	ciphertext, err := Encrypt(key, data)
	if err != nil {
		return err
	}

	// Prepend salt to ciphertext for storage.
	payload := append(salt, ciphertext...)

	if err := os.MkdirAll(filepath.Dir(v.Path), 0700); err != nil {
		return err
	}

	return os.WriteFile(v.Path, payload, 0600)
}

// Load reads the vault file from disk and decrypts it using the provided passphrase.
func (v *Vault) Load(passphrase string) (EnvMap, error) {
	if len(passphrase) == 0 {
		return nil, errors.New("passphrase must not be empty")
	}

	payload, err := os.ReadFile(v.Path)
	if err != nil {
		return nil, err
	}

	const saltLen = 32
	if len(payload) < saltLen {
		return nil, errors.New("vault file is corrupted or too short")
	}

	salt := payload[:saltLen]
	ciphertext := payload[saltLen:]

	key, _, err := DeriveKey(passphrase, salt)
	if err != nil {
		return nil, err
	}

	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}

	var env EnvMap
	if err := json.Unmarshal(plaintext, &env); err != nil {
		return nil, errors.New("failed to parse vault contents: wrong passphrase or corrupted data")
	}

	return env, nil
}

// Exists reports whether the vault file exists on disk.
func (v *Vault) Exists() bool {
	_, err := os.Stat(v.Path)
	return !os.IsNotExist(err)
}
