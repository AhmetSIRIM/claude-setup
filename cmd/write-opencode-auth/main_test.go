package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCredentials(t *testing.T) {
	credentialsPath := filepath.Join(t.TempDir(), "opencode", "auth.json")

	if err := writeCredentials(credentialsPath, "key-1"); err != nil {
		t.Fatalf("writeCredentials: %v", err)
	}

	content, err := os.ReadFile(credentialsPath)
	if err != nil {
		t.Fatalf("read credentials: %v", err)
	}
	var credentials map[string]map[string]string
	if err := json.Unmarshal(content, &credentials); err != nil {
		t.Fatalf("credentials are not JSON: %v (%s)", err, content)
	}
	if got := credentials["opencode-go"]; got["type"] != "api" || got["key"] != "key-1" {
		t.Errorf("opencode-go credential = %v, want type api and key key-1", got)
	}
	assertMode(t, credentialsPath, 0o600)
}

func TestWriteCredentialsNarrowsAnExistingFile(t *testing.T) {
	credentialsPath := filepath.Join(t.TempDir(), "auth.json")
	if err := os.WriteFile(credentialsPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	if err := writeCredentials(credentialsPath, "key-2"); err != nil {
		t.Fatalf("writeCredentials: %v", err)
	}

	assertMode(t, credentialsPath, 0o600)
}

func TestRunFailsWithoutKey(t *testing.T) {
	t.Setenv("OPENCODE_GO_KEY", "")
	if err := run(); err == nil {
		t.Fatal("run succeeded with an empty OPENCODE_GO_KEY")
	}
}

func assertMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Errorf("mode of %s = %o, want %o", path, got, want)
	}
}
