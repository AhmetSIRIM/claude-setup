// Command write-opencode-auth writes the opencode credentials file from the
// OPENCODE_GO_KEY environment variable, so the weekly digest can call the Go plan.
//
// The file is ~/.local/share/opencode/auth.json, readable by the owner only.
//
// Needs: OPENCODE_GO_KEY.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "write-opencode-auth:", err)
		os.Exit(1)
	}
}

func run() error {
	apiKey := os.Getenv("OPENCODE_GO_KEY")
	if apiKey == "" {
		return errors.New("OPENCODE_GO_KEY secret is empty")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if err := writeCredentials(filepath.Join(homeDir, ".local", "share", "opencode", "auth.json"), apiKey); err != nil {
		return err
	}
	fmt.Println("credentials written")
	return nil
}

// writeCredentials writes [apiKey] as the opencode-go API credential to
// [credentialsPath], creating its directory, and leaves the file at mode 0600 even
// when it already existed with wider permissions.
func writeCredentials(credentialsPath, apiKey string) error {
	content, err := json.Marshal(map[string]any{
		"opencode-go": map[string]string{"type": "api", "key": apiKey},
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(credentialsPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(credentialsPath, content, 0o600); err != nil {
		return err
	}
	return os.Chmod(credentialsPath, 0o600)
}
