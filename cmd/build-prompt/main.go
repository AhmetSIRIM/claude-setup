// Command build-prompt assembles the weekly digest prompt.
//
// The prompt gives the model two narrow tasks:
//  1. Newsletter: summarize ONLY the changelog sections newer than the last covered
//     version, parsed here from prior issue titles ("through vX.Y.Z"). The model never
//     decides what is new and never relates items to the setup.
//  2. Breakage check: report a problem ONLY for something the setup actually uses that
//     the current docs no longer support. Suggesting unused features is forbidden.
//
// It reads the setup from the working directory and the fetched inputs from /tmp,
// prints the prompt to stdout, and writes /tmp/digest-meta.env with LATEST_VERSION and
// NEW_SECTION_COUNT for the workflow. A missing input under /tmp counts as empty; any
// other read error stops the run.
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

//go:embed prompt.tmpl
var promptTemplateText string

var promptTemplate = template.Must(template.New("prompt").Parse(promptTemplateText))

const (
	setupFileCap = 25_000
	newsCap      = 35_000
	// On the first run no digest has covered a version yet.
	firstRunSectionLimit = 10
)

var (
	// The hook commands settings.template.json wires, read for the Claude Code
	// mechanisms they rely on (hook input fields, config file locations).
	hookSources = []string{"cmd/announce-git-mode/main.go", "cmd/session-start-doctor/main.go"}
	docNames    = []string{
		"llms.txt", "memory.md", "hooks.md", "agent-teams.md",
		"permission-modes.md", "auto-mode-config.md", "costs.md", "cross-session-messaging.md",
	}
	coveredVersionPattern  = regexp.MustCompile(`through v(\d+\.\d+\.\d+)`)
	changelogHeadingPrefix = regexp.MustCompile(`(?m)^## `)
	sectionVersionPattern  = regexp.MustCompile(`^(\d+\.\d+\.\d+)`)
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "build-prompt:", err)
		os.Exit(1)
	}
}

func run() error {
	const tmpDir = "/tmp"
	prompt, digestMeta, err := buildPrompt(".", tmpDir)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "digest-meta.env"), []byte(digestMeta), 0o644); err != nil {
		return err
	}
	_, err = fmt.Print(prompt)
	return err
}

// buildPrompt returns the prompt and the digest-meta.env content for the setup in
// [rootDir] and the fetched inputs in [tmpDir].
func buildPrompt(rootDir, tmpDir string) (prompt, digestMeta string, err error) {
	setupText, err := readSetup(rootDir)
	if err != nil {
		return "", "", err
	}
	prior, err := readOptional(filepath.Join(tmpDir, "prior-issues.md"))
	if err != nil {
		return "", "", err
	}
	changelog, err := readOptional(filepath.Join(tmpDir, "changelog.md"))
	if err != nil {
		return "", "", err
	}
	newSections, latestVersion := sectionsSinceLastDigest(changelog, lastCoveredVersion(prior))
	docs, err := readDocs(tmpDir)
	if err != nil {
		return "", "", err
	}

	if latestVersion == "" {
		latestVersion = "unknown"
	}
	digestMeta = fmt.Sprintf("LATEST_VERSION=%s\nNEW_SECTION_COUNT=%d\n", latestVersion, len(newSections))

	var promptBuilder strings.Builder
	err = promptTemplate.Execute(&promptBuilder, map[string]any{
		"NewSections": len(newSections) > 0,
		"News":        truncateRunes(strings.Join(newSections, "\n\n"), newsCap),
		"HasPrior":    strings.TrimSpace(prior) != "",
		"Prior":       prior,
		"Docs":        docs,
		"Setup":       setupText,
	})
	return promptBuilder.String(), digestMeta, err
}

// readSetup returns the setup files under review, each under a "### <path>" heading and
// cut at [setupFileCap] runes.
func readSetup(rootDir string) (string, error) {
	setupPaths := []string{"CLAUDE.template.md"}
	rulePaths, err := filepath.Glob(filepath.Join(rootDir, "rules", "*.md"))
	if err != nil {
		return "", err
	}
	for _, rulePath := range rulePaths {
		setupPaths = append(setupPaths, "rules/"+filepath.Base(rulePath))
	}
	setupPaths = append(setupPaths, hookSources...)
	skillPaths, err := filepath.Glob(filepath.Join(rootDir, "skills", "*", "SKILL.md"))
	if err != nil {
		return "", err
	}
	// Ordered by skill name, so "a" comes before "a-b" as a path-wise sort would put it.
	sort.Slice(skillPaths, func(i, j int) bool {
		return filepath.Base(filepath.Dir(skillPaths[i])) < filepath.Base(filepath.Dir(skillPaths[j]))
	})
	for _, skillPath := range skillPaths {
		setupPaths = append(setupPaths, "skills/"+filepath.Base(filepath.Dir(skillPath))+"/SKILL.md")
	}
	setupPaths = append(setupPaths, "settings.template.json", "README.md")

	sections := make([]string, 0, len(setupPaths))
	for _, setupPath := range setupPaths {
		content, err := readText(filepath.Join(rootDir, filepath.FromSlash(setupPath)))
		if err != nil {
			return "", err
		}
		if utf8RuneCount(content) > setupFileCap {
			content = truncateRunes(content, setupFileCap) + "\n[... truncated ...]"
		}
		sections = append(sections, "### "+setupPath+"\n"+content)
	}
	return strings.Join(sections, "\n\n"), nil
}

// readDocs returns the fetched docs pages whole, each under a "### <name>" heading;
// a page that was not fetched is left out.
func readDocs(tmpDir string) (string, error) {
	var sections []string
	for _, docName := range docNames {
		content, err := readText(filepath.Join(tmpDir, docName))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		sections = append(sections, "### "+docName+"\n"+content)
	}
	return strings.Join(sections, "\n\n"), nil
}

// lastCoveredVersion returns the highest "through vX.Y.Z" version in [prior], or ""
// when no digest has covered one.
func lastCoveredVersion(prior string) string {
	lastCovered := ""
	for _, match := range coveredVersionPattern.FindAllStringSubmatch(prior, -1) {
		if lastCovered == "" || compareVersions(match[1], lastCovered) > 0 {
			lastCovered = match[1]
		}
	}
	return lastCovered
}

// sectionsSinceLastDigest returns the changelog sections newer than [lastCovered],
// newest first, and the newest version in [changelog]. With no covered version it
// returns the newest [firstRunSectionLimit] sections.
func sectionsSinceLastDigest(changelog, lastCovered string) (newSections []string, latestVersion string) {
	for _, section := range changelogHeadingPrefix.Split(changelog, -1)[1:] {
		match := sectionVersionPattern.FindStringSubmatch(section)
		if match == nil {
			continue
		}
		version := match[1]
		if latestVersion == "" {
			latestVersion = version
		}
		if lastCovered != "" && compareVersions(version, lastCovered) <= 0 {
			break
		}
		newSections = append(newSections, "## "+strings.TrimSpace(section))
	}
	if lastCovered == "" && len(newSections) > firstRunSectionLimit {
		newSections = newSections[:firstRunSectionLimit]
	}
	return newSections, latestVersion
}

// compareVersions compares two "X.Y.Z" versions part by part as numbers. Both come
// from the digit-only version patterns above, so every part parses.
func compareVersions(left, right string) int {
	leftParts, rightParts := strings.Split(left, "."), strings.Split(right, ".")
	for i := range leftParts {
		leftNumber, _ := strconv.Atoi(leftParts[i])
		rightNumber, _ := strconv.Atoi(rightParts[i])
		if leftNumber != rightNumber {
			return leftNumber - rightNumber
		}
	}
	return 0
}

// readOptional returns the file's text, or "" when the file does not exist.
func readOptional(path string) (string, error) {
	content, err := readText(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	return content, err
}

// readText returns the file's content with invalid UTF-8 replaced by U+FFFD and
// CRLF line endings turned into LF.
func readText(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ToValidUTF8(string(content), "�"), "\r\n", "\n"), nil
}

func utf8RuneCount(text string) int { return len([]rune(text)) }

// truncateRunes returns at most [limit] runes of [text], so a cut never splits a
// multi-byte character.
func truncateRunes(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit])
}
