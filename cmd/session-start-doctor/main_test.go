package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const templateOnlyClaudeMarkdown = "<!-- fill this in -->\nAbout me:\n\n- Detailed working rules live in ~/.claude/rules/\n"

func TestDiagnoseHealthySetupIsSilent(t *testing.T) {
	configDir := t.TempDir()
	writeFile(t, configDir, "settings.json", `{"enabledPlugins": {"gopls-lsp@official": true}}`)
	writeFile(t, configDir, "plugins/installed_plugins.json", `{"plugins": {"gopls-lsp@official": []}}`)
	writeFile(t, configDir, "CLAUDE.md", templateOnlyClaudeMarkdown+"I build Android apps.\n")

	if warnings := diagnose(configDir, everyBinaryExists); len(warnings) != 0 {
		t.Errorf("warnings = %q, want none", warnings)
	}
}

func TestDiagnoseReportsEachProblemInSettingsOrder(t *testing.T) {
	configDir := t.TempDir()
	toolsDir := filepath.Join(configDir, "tools")
	writeFile(t, toolsDir, "package.json", `{"dependencies": {"ccstatusline": "^2.0.0"}}`)
	writeFile(t, toolsDir, "node_modules/ccstatusline/package.json", `{"version": "1.9.0"}`)
	writeFile(t, toolsDir, "node_modules/.bin/ccstatusline", "")
	writeFile(t, configDir, "settings.json", `{
	  "enabledPlugins": {"swift-lsp@official": true, "gopls-lsp@official": true, "jdtls-lsp@official": false},
	  "hooks": {"SessionStart": [{"hooks": [{"type": "command", "command": "`+configDir+`/missing-hook --flag"}]}]},
	  "statusLine": {"type": "command", "command": "\"`+toolsDir+`/node_modules/.bin/ccstatusline\" --config \"`+configDir+`/missing.json\""}
	}`)
	writeFile(t, configDir, "plugins/installed_plugins.json", `{"plugins": {"gopls-lsp@official": []}}`)
	writeFile(t, configDir, "CLAUDE.md", templateOnlyClaudeMarkdown)
	if err := os.Symlink(filepath.Join(configDir, "gone"), filepath.Join(configDir, "rules")); err != nil {
		t.Fatal(err)
	}

	warnings := diagnose(configDir, func(name string) bool { return name != "sourcekit-lsp" })

	want := []string{
		"plugin 'swift-lsp' is enabled but its language server 'sourcekit-lsp' is not installed",
		"plugin 'swift-lsp@official' is enabled in settings but not installed; run: claude plugin install swift-lsp@official",
		"~/.claude/rules is a broken symlink",
		"~/.claude/CLAUDE.md is still the unfilled template; make it yours",
		"settings references a missing hook: " + configDir + "/missing-hook --flag",
		"ccstatusline 1.9.0 is installed but ^2.0.0 is pinned; run: npm ci --prefix " + toolsDir,
		`statusLine references a missing file: "` + configDir + `/missing.json"`,
	}
	if !reflect.DeepEqual(warnings, want) {
		t.Errorf("warnings =\n%q\nwant\n%q", warnings, want)
	}
}

func TestDiagnoseWithoutUsableSettings(t *testing.T) {
	configDir := t.TempDir()
	if warnings := diagnose(configDir, everyBinaryExists); len(warnings) != 1 {
		t.Fatalf("warnings = %q, want one unreadable-settings warning", warnings)
	}

	writeFile(t, configDir, "settings.json", "{not json")
	if warnings := diagnose(configDir, everyBinaryExists); len(warnings) != 1 {
		t.Fatalf("warnings = %q, want one unreadable-settings warning", warnings)
	}
}

func TestDiagnoseMissingClaudeMarkdown(t *testing.T) {
	configDir := t.TempDir()
	writeFile(t, configDir, "settings.json", `{}`)

	want := []string{"~/.claude/CLAUDE.md is missing; copy CLAUDE.template.md from the setup repo and fill it in"}
	if warnings := diagnose(configDir, everyBinaryExists); !reflect.DeepEqual(warnings, want) {
		t.Errorf("warnings = %q, want %q", warnings, want)
	}
}

func everyBinaryExists(string) bool { return true }

func writeFile(t *testing.T, dir, relativePath, content string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
