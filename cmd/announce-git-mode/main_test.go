package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestAnnouncement(t *testing.T) {
	testCases := []struct {
		name         string
		previousMode string
		seenBefore   bool
		currentMode  string
		wantAnnounce bool
		wantPrefix   string
		wantMention  string
	}{
		{name: "first prompt in manual mode names git mode ask", currentMode: "default", wantAnnounce: true, wantPrefix: "Git mode: ask", wantMention: "git.md"},
		{name: "first prompt in auto mode names git mode ask", currentMode: "auto", wantAnnounce: true, wantPrefix: "Git mode: ask", wantMention: "git.md"},
		{name: "first prompt with bypass announces the autonomous session", currentMode: "bypassPermissions", wantAnnounce: true, wantPrefix: "Autonomous session", wantMention: "Git mode: free"},
		{name: "a missing mode is named, not hidden", currentMode: "", wantAnnounce: true, wantPrefix: "Git mode: ask", wantMention: "carried no permission mode"},
		{name: "an unchanged mode stays silent", previousMode: "auto", seenBefore: true, currentMode: "auto", wantAnnounce: false},
		{name: "a switch into bypass is announced with the change", previousMode: "auto", seenBefore: true, currentMode: "bypassPermissions", wantAnnounce: true, wantPrefix: `Permission mode is now "bypassPermissions". Autonomous session`, wantMention: "Git mode: free"},
		{name: "a switch out of bypass returns to ask", previousMode: "bypassPermissions", seenBefore: true, currentMode: "default", wantAnnounce: true, wantPrefix: `Permission mode is now "default". Git mode: ask`, wantMention: "git.md"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			line, announce := announcement(testCase.previousMode, testCase.seenBefore, testCase.currentMode)
			if announce != testCase.wantAnnounce {
				t.Fatalf("announce = %v, want %v (line %q)", announce, testCase.wantAnnounce, line)
			}
			if !announce {
				return
			}
			if !strings.HasPrefix(line, testCase.wantPrefix) {
				t.Errorf("line = %q, want prefix %q", line, testCase.wantPrefix)
			}
			if !strings.Contains(line, testCase.wantMention) {
				t.Errorf("line = %q, want it to mention %q", line, testCase.wantMention)
			}
		})
	}
}

func TestStateRoundTrip(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "announce-git-mode", "session-1")

	previousMode, seenBefore, err := readState(statePath)
	if err != nil || seenBefore || previousMode != "" {
		t.Fatalf("fresh state = (%q, %v, %v), want (\"\", false, nil)", previousMode, seenBefore, err)
	}
	if err := writeState(statePath, "bypassPermissions"); err != nil {
		t.Fatalf("writeState: %v", err)
	}
	previousMode, seenBefore, err = readState(statePath)
	if err != nil || !seenBefore || previousMode != "bypassPermissions" {
		t.Fatalf("state after write = (%q, %v, %v), want (\"bypassPermissions\", true, nil)", previousMode, seenBefore, err)
	}
}
