package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLastCoveredVersion(t *testing.T) {
	testCases := []struct {
		name  string
		prior string
		want  string
	}{
		{name: "no digest yet", prior: "", want: ""},
		{name: "the highest version wins, compared as numbers", prior: "digest (through v2.1.9)\ndigest (through v2.1.10)\ndigest (through v2.0.99)", want: "2.1.10"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := lastCoveredVersion(testCase.prior); got != testCase.want {
				t.Errorf("lastCoveredVersion = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestSectionsSinceLastDigest(t *testing.T) {
	changelog := "# Changelog\n\n## 2.1.12\n\n- c\n\n## Unversioned notes\n\n## 2.1.11\n\n- b\n\n## 2.1.10\n\n- a\n"

	newSections, latestVersion := sectionsSinceLastDigest(changelog, "2.1.10")
	if want := []string{"## 2.1.12\n\n- c", "## 2.1.11\n\n- b"}; !reflect.DeepEqual(newSections, want) {
		t.Errorf("sections since 2.1.10 = %q, want %q", newSections, want)
	}
	if latestVersion != "2.1.12" {
		t.Errorf("latestVersion = %q, want 2.1.12", latestVersion)
	}

	var manyReleases strings.Builder
	manyReleases.WriteString("# Changelog\n")
	for patch := 20; patch > 0; patch-- {
		fmt.Fprintf(&manyReleases, "\n## 2.1.%d\n\n- item\n", patch)
	}
	firstRunSections, _ := sectionsSinceLastDigest(manyReleases.String(), "")
	if len(firstRunSections) != firstRunSectionLimit || !strings.HasPrefix(firstRunSections[0], "## 2.1.20") {
		t.Errorf("first run returned %d sections starting %q, want %d starting with 2.1.20", len(firstRunSections), firstRunSections[0], firstRunSectionLimit)
	}
}

func TestBuildPrompt(t *testing.T) {
	rootDir, tmpDir := t.TempDir(), t.TempDir()
	writeFile(t, rootDir, "CLAUDE.template.md", "template")
	writeFile(t, rootDir, "rules/b.md", "rule b")
	writeFile(t, rootDir, "rules/a.md", strings.Repeat("ü", setupFileCap+5))
	writeFile(t, rootDir, "hooks/doctor.sh", "hook")
	writeFile(t, rootDir, "skills/a-b/SKILL.md", "skill a-b")
	writeFile(t, rootDir, "skills/a/SKILL.md", "skill a")
	writeFile(t, rootDir, "settings.template.json", "{}")
	writeFile(t, rootDir, "README.md", "readme")
	writeFile(t, tmpDir, "changelog.md", "# Changelog\n\n## 2.1.2\n\n- new\n\n## 2.1.1\n\n- old\n")
	writeFile(t, tmpDir, "hooks.md", "hooks doc")
	longDoc := strings.Repeat("d", setupFileCap+5)
	writeFile(t, tmpDir, "memory.md", longDoc)
	writeFile(t, tmpDir, "prior-issues.md", "Claude Code weekly digest (through v2.1.1)\r\nclosing comment\r\n")

	prompt, digestMeta, err := buildPrompt(rootDir, tmpDir)
	if err != nil {
		t.Fatalf("buildPrompt: %v", err)
	}

	if want := "LATEST_VERSION=2.1.2\nNEW_SECTION_COUNT=1\n"; digestMeta != want {
		t.Errorf("digestMeta = %q, want %q", digestMeta, want)
	}
	for _, want := range []string{
		"## 2.1.2\n\n- new",
		"## Previously reported\n",
		"### hooks.md\nhooks doc",
		"### memory.md\n" + longDoc + "\n\n",
		"closing comment\n",
		strings.Repeat("ü", setupFileCap) + "\n[... truncated ...]",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt is missing %q", truncateRunes(want, 60))
		}
	}
	if strings.Contains(prompt, "## 2.1.1\n\n- old") {
		t.Error("prompt carries a section the last digest already covered")
	}
	order := []string{"### CLAUDE.template.md", "### rules/a.md", "### rules/b.md", "### hooks/doctor.sh", "### skills/a/SKILL.md", "### skills/a-b/SKILL.md", "### settings.template.json", "### README.md"}
	lastIndex := -1
	for _, heading := range order {
		index := strings.Index(prompt, heading)
		if index <= lastIndex {
			t.Fatalf("%q is out of order in the setup section", heading)
		}
		lastIndex = index
	}
}

func TestBuildPromptWithoutPriorOrReleases(t *testing.T) {
	rootDir := t.TempDir()
	for _, name := range []string{"CLAUDE.template.md", "settings.template.json", "README.md"} {
		writeFile(t, rootDir, name, name)
	}

	prompt, digestMeta, err := buildPrompt(rootDir, t.TempDir())
	if err != nil {
		t.Fatalf("buildPrompt: %v", err)
	}

	if want := "LATEST_VERSION=unknown\nNEW_SECTION_COUNT=0\n"; digestMeta != want {
		t.Errorf("digestMeta = %q, want %q", digestMeta, want)
	}
	if !strings.Contains(prompt, "No new releases since the last digest.") {
		t.Error("prompt lacks the empty newsletter task")
	}
	if strings.Contains(prompt, "Previously reported") {
		t.Error("prompt has a previously-reported section without prior issues")
	}
}

func TestBuildPromptFailsWithoutASetupFile(t *testing.T) {
	if _, _, err := buildPrompt(t.TempDir(), t.TempDir()); err == nil {
		t.Fatal("buildPrompt succeeded with no setup files, want an error")
	}
}

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
