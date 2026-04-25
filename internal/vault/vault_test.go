package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVaultSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	v := New("test", dir)

	env := EnvMap{
		"API_KEY":  "super-secret-123",
		"DB_URL":   "postgres://localhost/mydb",
		"APP_PORT": "8080",
	}
	passphrase := "correct-horse-battery-staple"

	if err := v.Save(env, passphrase); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := v.Load(passphrase)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	for k, want := range env {
		if got := loaded[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestVaultLoadWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	v := New("test", dir)

	env := EnvMap{"SECRET": "value"}

	if err := v.Save(env, "correct-passphrase"); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	_, err := v.Load("wrong-passphrase")
	if err == nil {
		t.Fatal("Load() with wrong passphrase should return an error")
	}
}

func TestVaultExists(t *testing.T) {
	dir := t.TempDir()
	v := New("myenv", dir)

	if v.Exists() {
		t.Fatal("Exists() should be false before saving")
	}

	if err := v.Save(EnvMap{"K": "V"}, "passphrase"); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if !v.Exists() {
		t.Fatal("Exists() should be true after saving")
	}
}

func TestVaultSaveEmptyPassphrase(t *testing.T) {
	dir := t.TempDir()
	v := New("test", dir)

	err := v.Save(EnvMap{"K": "V"}, "")
	if err == nil {
		t.Fatal("Save() with empty passphrase should return an error")
	}
}

func TestVaultFilePermissions(t *testing.T) {
	dir := t.TempDir()
	v := New("secure", dir)

	if err := v.Save(EnvMap{"X": "Y"}, "passphrase"); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, "secure.vault"))
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file permissions 0600, got %04o", perm)
	}
}
