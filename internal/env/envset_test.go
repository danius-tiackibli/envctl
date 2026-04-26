package env

import (
	"testing"
)

func TestNewEnvSet(t *testing.T) {
	es, err := New("production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Name != "production" {
		t.Errorf("expected name 'production', got %q", es.Name)
	}
	if len(es.Variables) != 0 {
		t.Errorf("expected empty variables map")
	}
}

func TestNewEnvSetEmptyName(t *testing.T) {
	_, err := New("  ")
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestEnvSetSetAndGet(t *testing.T) {
	es, _ := New("test")
	if err := es.Set("DB_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := es.Get("DB_HOST")
	if !ok {
		t.Fatal("expected key DB_HOST to exist")
	}
	if v != "localhost" {
		t.Errorf("expected 'localhost', got %q", v)
	}
}

func TestEnvSetInvalidKey(t *testing.T) {
	es, _ := New("test")
	if err := es.Set("", "value"); err == nil {
		t.Error("expected error for empty key")
	}
	if err := es.Set("INVALID KEY", "value"); err == nil {
		t.Error("expected error for key with space")
	}
	if err := es.Set("INVALID=KEY", "value"); err == nil {
		t.Error("expected error for key with '='")
	}
}

func TestEnvSetDelete(t *testing.T) {
	es, _ := New("test")
	_ = es.Set("TO_DELETE", "bye")
	es.Delete("TO_DELETE")
	if _, ok := es.Get("TO_DELETE"); ok {
		t.Error("expected key to be deleted")
	}
}

func TestEnvSetToExportLines(t *testing.T) {
	es, _ := New("test")
	_ = es.Set("Z_VAR", "z")
	_ = es.Set("A_VAR", "a")
	lines := es.ToExportLines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "export A_VAR=a" {
		t.Errorf("expected sorted first line 'export A_VAR=a', got %q", lines[0])
	}
}

func TestEnvSetMerge(t *testing.T) {
	base, _ := New("base")
	_ = base.Set("KEY1", "original")
	_ = base.Set("KEY2", "keep")

	override, _ := New("override")
	_ = override.Set("KEY1", "overridden")
	_ = override.Set("KEY3", "new")

	base.Merge(override)

	if v, _ := base.Get("KEY1"); v != "overridden" {
		t.Errorf("expected KEY1 to be overridden, got %q", v)
	}
	if v, _ := base.Get("KEY2"); v != "keep" {
		t.Errorf("expected KEY2 to remain 'keep', got %q", v)
	}
	if _, ok := base.Get("KEY3"); !ok {
		t.Error("expected KEY3 to exist after merge")
	}
}
