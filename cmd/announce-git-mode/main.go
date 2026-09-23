// Command announce-git-mode is the UserPromptSubmit hook that names the session's
// git mode.
//
// It reads the hook input from stdin and prints one line that Claude Code adds to
// the session context. The line is printed when a session's permission mode is
// first seen and again whenever it changes, so a switch into bypass permissions
// during the session is announced on the next prompt. A session in
// bypassPermissions mode is an autonomous run: the line says so, points at the
// rule that defines it, and names git mode free. Every other mode gets git mode
// ask. The same line is the user's warning, shown wherever the transcript shows
// hook output.
//
// The prompt hook is used because its input carries permission_mode; the
// SessionStart input does not. The last announced mode is kept per session in a
// file under the OS temp dir.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const bypassPermissionsMode = "bypassPermissions"

type hookInput struct {
	SessionID      string `json:"session_id"`
	PermissionMode string `json:"permission_mode"`
}

func main() {
	rawInput, err := io.ReadAll(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "announce-git-mode: read hook input: %v\n", err)
		os.Exit(1)
	}
	var input hookInput
	if err := json.Unmarshal(rawInput, &input); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "announce-git-mode: hook input is not JSON: %v\n", err)
		os.Exit(1)
	}

	statePath := filepath.Join(os.TempDir(), "announce-git-mode", input.SessionID)
	previousMode, seenBefore, stateErr := readState(statePath)
	line, announce := announcement(previousMode, seenBefore, input.PermissionMode)
	if stateErr == nil {
		stateErr = writeState(statePath, input.PermissionMode)
	}
	if stateErr != nil {
		// The line still goes out; the note makes the repeat on every prompt
		// explain itself instead of looking like a stuck hook.
		line += fmt.Sprintf(" (announce-git-mode cannot keep its state, so this line repeats: %v)", stateErr)
		announce = true
	}
	if announce {
		fmt.Println(line)
	}
}

// announcement is the context line for a session whose permission mode is now
// [currentMode], and whether it should be printed at all. [previousMode] is the
// mode announced last in this session; [seenBefore] is false on the session's first
// prompt. An empty current mode means the hook input carried none; the line then
// says so instead of passing the default off as a reading.
func announcement(previousMode string, seenBefore bool, currentMode string) (line string, announce bool) {
	if seenBefore && previousMode == currentMode {
		return "", false
	}
	switch currentMode {
	case bypassPermissionsMode:
		line = "Autonomous session: bypass permissions is on, so no rule waits for the " +
			"user; every ask becomes a decision and a line in the final report. See " +
			"~/.claude/rules/autonomous-session.md. Git mode: free. Say 'ask' or 'plan' " +
			"to switch for this session."
	case "":
		line = "Git mode: ask (default; the hook input carried no permission mode). Say " +
			"'plan' or 'free' to switch for this session; see ~/.claude/rules/git.md for " +
			"what each mode covers."
	default:
		line = "Git mode: ask (default). Say 'plan' or 'free' to switch for this session; " +
			"see ~/.claude/rules/git.md for what each mode covers."
	}
	if seenBefore {
		line = fmt.Sprintf("Permission mode is now %q. ", currentMode) + line
	}
	return line, true
}

// readState returns the mode last announced for the session, and whether one was
// recorded. A missing file is the first prompt, not an error.
func readState(statePath string) (previousMode string, seenBefore bool, err error) {
	content, err := os.ReadFile(statePath)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(content), true, nil
}

func writeState(statePath, mode string) error {
	if err := os.MkdirAll(filepath.Dir(statePath), 0o700); err != nil {
		return err
	}
	return os.WriteFile(statePath, []byte(mode), 0o600)
}
