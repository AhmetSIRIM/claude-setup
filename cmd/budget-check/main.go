// Command budget-check reads the OpenCode Go plan usage windows and decides whether
// the weekly digest can run.
//
// The plan enforces three windows (rolling, weekly, monthly); a request is refused
// while any of them is rate-limited. The outcome goes to GITHUB_OUTPUT as
// skip=true|false; a limited window also leaves a line in GITHUB_STEP_SUMMARY.
//
// Needs: OPENCODE_GO_KEY, GITHUB_OUTPUT. Optional: GITHUB_STEP_SUMMARY.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	goPlanAPIBaseURL = "https://opencode.ai/zen/go/v1"
	requestTimeout   = 30 * time.Second
	// The usage endpoint reports a window the plan is refusing requests for with
	// this status; the chat endpoint answers 429 for the same condition.
	rateLimitedStatus = "rate-limited"
)

// UsageWindow is one usage window as the API reports it.
type UsageWindow struct {
	Status      string  `json:"status"`
	UsedPercent float64 `json:"percent"`
	ResetsAt    string  `json:"resetsAt"`
}

type usageResponse struct {
	Usage map[string]UsageWindow `json:"usage"`
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "budget-check:", err)
		os.Exit(1)
	}
}

func run() error {
	apiKey, err := requiredEnv("OPENCODE_GO_KEY")
	if err != nil {
		return err
	}

	githubOutputPath, err := requiredEnv("GITHUB_OUTPUT")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	usageWindows, err := fetchUsageWindows(ctx, http.DefaultClient, goPlanAPIBaseURL, apiKey)
	if err != nil {
		return err
	}

	for _, windowName := range sortedWindowNames(usageWindows) {
		window := usageWindows[windowName]
		fmt.Printf("%s: %s, %.0f%% used, resets at %s\n", windowName, window.Status, window.UsedPercent, window.ResetsAt)
	}
	limitedWindowDescriptions := describeLimitedWindows(usageWindows)
	skipDigest := len(limitedWindowDescriptions) > 0
	if err := appendLineToFile(githubOutputPath, fmt.Sprintf("skip=%t", skipDigest)); err != nil {
		return err
	}
	if !skipDigest {
		return nil
	}
	skipReason := "Go plan " + strings.Join(limitedWindowDescriptions, "; ")
	fmt.Printf("::notice::doc-drift-check skipped: %s\n", skipReason)
	if summaryPath := os.Getenv("GITHUB_STEP_SUMMARY"); summaryPath != "" {
		return appendLineToFile(summaryPath, "Skipped: "+skipReason)
	}
	return nil
}

// fetchUsageWindows returns the plan's usage windows as the API reports them. A
// non-200 answer or a body without windows is returned as an error carrying the
// status and body, so the run log names the cause.
func fetchUsageWindows(ctx context.Context, client *http.Client, apiBaseURL, apiKey string) (map[string]UsageWindow, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBaseURL+"/usage", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("usage request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 64<<10))
	if err != nil {
		return nil, fmt.Errorf("usage response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("usage request failed: HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var decodedResponse usageResponse
	if err := json.Unmarshal(responseBody, &decodedResponse); err != nil || len(decodedResponse.Usage) == 0 {
		return nil, fmt.Errorf("usage response has no windows: %s", strings.TrimSpace(string(responseBody)))
	}
	return decodedResponse.Usage, nil
}

// describeLimitedWindows names the rate-limited windows with their reset times, in a
// stable order. An empty result means the digest may run.
func describeLimitedWindows(usageWindows map[string]UsageWindow) []string {
	var descriptions []string
	for _, windowName := range sortedWindowNames(usageWindows) {
		if window := usageWindows[windowName]; window.Status == rateLimitedStatus {
			descriptions = append(descriptions, fmt.Sprintf("%s window at %.0f%%, resets at %s", windowName, window.UsedPercent, window.ResetsAt))
		}
	}
	return descriptions
}

func sortedWindowNames(usageWindows map[string]UsageWindow) []string {
	windowNames := make([]string, 0, len(usageWindows))
	for windowName := range usageWindows {
		windowNames = append(windowNames, windowName)
	}
	sort.Strings(windowNames)
	return windowNames
}

func requiredEnv(variableName string) (string, error) {
	value := os.Getenv(variableName)
	if value == "" {
		return "", errors.New(variableName + " is not set")
	}
	return value, nil
}

func appendLineToFile(filePath, line string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = fmt.Fprintln(file, line)
	return err
}
