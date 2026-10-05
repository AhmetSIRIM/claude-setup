// Command session-start-doctor is the SessionStart hook that checks this setup's
// health. It prints one "setup-doctor:" line per problem and stays silent when the
// setup is healthy, so a healthy session gets no extra context.
//
// It catches configuration that breaks without an error: a plugin enabled in settings
// but not installed, an enabled code-intelligence plugin whose language server is
// missing, a hook or status line file that settings reference but that is absent, a
// broken ~/.claude symlink, a missing or unfilled ~/.claude/CLAUDE.md, and a pinned
// npm tool installed at another version than its package.json names.
//
// CLAUDE_DIR points the check at another config directory. The command always exits 0:
// a SessionStart hook reports through its output, not its status.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// languageServerBinaries maps a code-intelligence plugin to the binary it starts.
var languageServerBinaries = map[string]string{
	"kotlin-lsp": "kotlin-language-server", "swift-lsp": "sourcekit-lsp",
	"gopls-lsp": "gopls", "typescript-lsp": "typescript-language-server",
	"jdtls-lsp": "jdtls", "pyright-lsp": "pyright-langserver",
	"rust-analyzer-lsp": "rust-analyzer", "clangd-lsp": "clangd",
	"csharp-lsp": "csharp-ls", "lua-lsp": "lua-language-server", "php-lsp": "intelephense",
}

var (
	htmlCommentPattern = regexp.MustCompile(`(?s)<!--.*?-->`)
	// A line of the CLAUDE.md template that is not the owner's own writing: blank, a
	// "Heading:" label, or the pointer to the rules.
	templateScaffoldLine = regexp.MustCompile(`^\s*$|^[A-Z][A-Za-z &]+:\s*$|^- Detailed working rules `)
)

func main() {
	for _, warning := range diagnose(claudeDir(), lookupBinary) {
		fmt.Println("setup-doctor:", warning)
	}
}

func claudeDir() string {
	if dir := os.Getenv("CLAUDE_DIR"); dir != "" {
		return dir
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~/.claude"
	}
	return filepath.Join(homeDir, ".claude")
}

// diagnose returns the problems found in the config directory [configDir], using
// [binaryExists] to look for language servers.
func diagnose(configDir string, binaryExists func(string) bool) []string {
	settingsContent, err := os.ReadFile(filepath.Join(configDir, "settings.json"))
	if err != nil {
		return []string{"settings.json unreadable: " + err.Error()}
	}
	var settings struct {
		EnabledPlugins json.RawMessage `json:"enabledPlugins"`
		Hooks          json.RawMessage `json:"hooks"`
		StatusLine     struct {
			Type    string `json:"type"`
			Command string `json:"command"`
		} `json:"statusLine"`
	}
	if err := json.Unmarshal(settingsContent, &settings); err != nil {
		return []string{"settings.json unreadable: " + err.Error()}
	}
	enabledPluginIDs, err := enabledPlugins(settings.EnabledPlugins)
	if err != nil {
		return []string{"settings.json unreadable: enabledPlugins: " + err.Error()}
	}

	var warnings []string
	for _, pluginID := range enabledPluginIDs {
		pluginName, _, _ := strings.Cut(pluginID, "@")
		if binary, known := languageServerBinaries[pluginName]; known && !binaryExists(binary) {
			warnings = append(warnings, fmt.Sprintf("plugin '%s' is enabled but its language server '%s' is not installed", pluginName, binary))
		}
	}
	warnings = append(warnings, uninstalledPluginWarnings(configDir, enabledPluginIDs)...)

	for _, entry := range []string{"rules", "skills", "CLAUDE.md"} {
		if isBrokenSymlink(filepath.Join(configDir, entry)) {
			warnings = append(warnings, fmt.Sprintf("~/.claude/%s is a broken symlink", entry))
		}
	}
	warnings = append(warnings, claudeMarkdownWarnings(filepath.Join(configDir, "CLAUDE.md"))...)

	hookCommands, err := hookCommands(settings.Hooks)
	if err != nil {
		warnings = append(warnings, "could not read the hooks in settings.json: "+err.Error())
	}
	for _, command := range hookCommands {
		firstWord, _, _ := strings.Cut(command, " ")
		executable := expandHome(firstWord)
		if (strings.HasPrefix(executable, "/") || strings.HasPrefix(command, "~")) && !pathExists(executable) {
			warnings = append(warnings, "settings references a missing hook: "+command)
		}
	}

	if settings.StatusLine.Type == "command" {
		warnings = append(warnings, statusLineWarnings(settings.StatusLine.Command)...)
	}
	return warnings
}

// enabledPlugins returns the IDs set to true in the enabledPlugins object, in the
// order settings.json lists them.
func enabledPlugins(rawPlugins json.RawMessage) ([]string, error) {
	pluginIDs, values, err := orderedObject(rawPlugins)
	if err != nil {
		return nil, err
	}
	var enabledIDs []string
	for _, pluginID := range pluginIDs {
		var enabled bool
		if json.Unmarshal(values[pluginID], &enabled) == nil && enabled {
			enabledIDs = append(enabledIDs, pluginID)
		}
	}
	return enabledIDs, nil
}

// uninstalledPluginWarnings names each enabled plugin that installed_plugins.json
// does not list; a plugin enabled but never installed does nothing. A missing file
// means nothing is installed through the marketplace yet, so it gives no warning.
func uninstalledPluginWarnings(configDir string, enabledPluginIDs []string) []string {
	content, err := os.ReadFile(filepath.Join(configDir, "plugins", "installed_plugins.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	var installed struct {
		Plugins map[string]json.RawMessage `json:"plugins"`
	}
	if err == nil {
		err = json.Unmarshal(content, &installed)
	}
	if err != nil {
		return []string{"could not read installed_plugins.json: " + err.Error()}
	}
	var warnings []string
	for _, pluginID := range enabledPluginIDs {
		if _, isInstalled := installed.Plugins[pluginID]; !isInstalled {
			warnings = append(warnings, fmt.Sprintf("plugin '%s' is enabled in settings but not installed; run: claude plugin install %s", pluginID, pluginID))
		}
	}
	return warnings
}

// claudeMarkdownWarnings reports a missing CLAUDE.md, or one that still holds only the
// template: HTML comments and the template's scaffold lines are dropped, and whatever
// remains is the owner's own writing.
func claudeMarkdownWarnings(claudeMarkdownPath string) []string {
	content, err := os.ReadFile(claudeMarkdownPath)
	if errors.Is(err, fs.ErrNotExist) {
		return []string{"~/.claude/CLAUDE.md is missing; copy CLAUDE.template.md from the setup repo and fill it in"}
	}
	if err != nil {
		return []string{"could not read ~/.claude/CLAUDE.md: " + err.Error()}
	}
	body := htmlCommentPattern.ReplaceAllString(string(content), "")
	for line := range strings.SplitSeq(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		if !templateScaffoldLine.MatchString(line) {
			return nil
		}
	}
	return []string{"~/.claude/CLAUDE.md is still the unfilled template; make it yours"}
}

// hookCommands returns every hook command in the hooks object, in the order
// settings.json lists the events.
func hookCommands(rawHooks json.RawMessage) ([]string, error) {
	eventNames, events, err := orderedObject(rawHooks)
	if err != nil {
		return nil, err
	}
	var commands []string
	for _, eventName := range eventNames {
		var matchers []struct {
			Hooks []struct {
				Command string `json:"command"`
			} `json:"hooks"`
		}
		if err := json.Unmarshal(events[eventName], &matchers); err != nil {
			return commands, fmt.Errorf("%s: %w", eventName, err)
		}
		for _, matcher := range matchers {
			for _, hook := range matcher.Hooks {
				if hook.Command != "" {
					commands = append(commands, hook.Command)
				}
			}
		}
	}
	return commands, nil
}

// statusLineWarnings checks every token of the status line [command] that expands to
// an absolute path, since the command may be a bare program with path arguments. For
// an npm tool it also compares the installed version with the package.json pin: a tool
// installed from an older lockfile keeps running, so the pin only holds when checked.
func statusLineWarnings(command string) []string {
	var warnings []string
	for token := range strings.FieldsSeq(command) {
		expanded := expandHome(expandEnv(strings.Trim(token, `"'`)))
		if strings.HasPrefix(expanded, "/") && !pathExists(expanded) {
			warnings = append(warnings, "statusLine references a missing file: "+token)
			continue
		}
		toolsDir, toolName, isNpmTool := strings.Cut(expanded, "/node_modules/.bin/")
		if !isNpmTool {
			continue
		}
		pinnedVersion, installedVersion, err := npmToolVersions(toolsDir, toolName)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("could not compare the pinned %s version: %v", toolName, err))
			continue
		}
		if pinnedVersion != "" && installedVersion != "" && strings.TrimLeft(pinnedVersion, "^~=") != installedVersion {
			warnings = append(warnings, fmt.Sprintf("%s %s is installed but %s is pinned; run: npm ci --prefix %s", toolName, installedVersion, pinnedVersion, toolsDir))
		}
	}
	return warnings
}

// npmToolVersions returns the version [toolsDir]/package.json pins for [toolName] and
// the version installed under its node_modules.
func npmToolVersions(toolsDir, toolName string) (pinnedVersion, installedVersion string, err error) {
	var manifest struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := readJSON(filepath.Join(toolsDir, "package.json"), &manifest); err != nil {
		return "", "", err
	}
	var installedManifest struct {
		Version string `json:"version"`
	}
	if err := readJSON(filepath.Join(toolsDir, "node_modules", toolName, "package.json"), &installedManifest); err != nil {
		return "", "", err
	}
	return manifest.Dependencies[toolName], installedManifest.Version, nil
}

// orderedObject returns the keys of the JSON object [raw] in document order, with
// their raw values. A null or absent object has no keys.
func orderedObject(raw json.RawMessage) ([]string, map[string]json.RawMessage, error) {
	values := map[string]json.RawMessage{}
	if len(bytes.TrimSpace(raw)) == 0 || string(bytes.TrimSpace(raw)) == "null" {
		return nil, values, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	if token, err := decoder.Token(); err != nil || token != json.Delim('{') {
		return nil, nil, fmt.Errorf("not an object")
	}
	var keys []string
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return nil, nil, err
		}
		key, _ := token.(string)
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return nil, nil, err
		}
		if _, seen := values[key]; !seen {
			keys = append(keys, key)
		}
		values[key] = value
	}
	return keys, values, nil
}

// lookupBinary reports whether [name] is on the PATH or in a usual install directory;
// a hook may run without the login shell's PATH. sourcekit-lsp ships inside Xcode and
// is found through xcrun.
func lookupBinary(name string) bool {
	if _, err := exec.LookPath(name); err == nil {
		return true
	}
	for _, dir := range []string{expandHome("~/go/bin"), "/opt/homebrew/bin", "/usr/local/bin"} {
		if info, err := os.Stat(filepath.Join(dir, name)); err == nil && info.Mode()&0o111 != 0 {
			return true
		}
	}
	if name == "sourcekit-lsp" {
		return exec.Command("xcrun", "--find", "sourcekit-lsp").Run() == nil
	}
	return false
}

func isBrokenSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0 && !pathExists(path)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readJSON(path string, target any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, target)
}

// expandHome replaces a leading "~" with the home directory.
func expandHome(path string) string {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return homeDir + path[1:]
}

// expandEnv replaces $NAME and ${NAME} with the variable's value, and keeps a
// reference to an unset variable as ${NAME} instead of dropping it.
func expandEnv(text string) string {
	return os.Expand(text, func(name string) string {
		if value, isSet := os.LookupEnv(name); isSet {
			return value
		}
		return "${" + name + "}"
	})
}
