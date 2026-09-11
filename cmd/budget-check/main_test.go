package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const threeWindowsBody = `{"usage":{
  "rolling":{"status":"ok","percent":12,"resetsAt":"2026-01-01T05:00:00.000Z"},
  "weekly":{"status":"ok","percent":34,"resetsAt":"2026-01-06T00:00:00.000Z"},
  "monthly":{"status":"rate-limited","percent":100,"resetsAt":"2026-02-01T00:00:00.000Z"}}}`

func TestFetchUsageWindows(t *testing.T) {
	testCases := []struct {
		name            string
		responseStatus  int
		responseBody    string
		wantWindowCount int
		wantErr         string
	}{
		{
			name:            "a 200 with windows is returned as the API sent it",
			responseStatus:  http.StatusOK,
			responseBody:    threeWindowsBody,
			wantWindowCount: 3,
		},
		{
			name:           "bad key is a real failure",
			responseStatus: http.StatusUnauthorized,
			responseBody:   `{"type":"error","error":{"type":"AuthError","message":"Unauthorized"}}`,
			wantErr:        "HTTP 401",
		},
		{
			name:           "no Go subscription is a real failure",
			responseStatus: http.StatusForbidden,
			responseBody:   `{"type":"error","error":{"type":"EntitlementError","message":"OpenCode Go subscription required."}}`,
			wantErr:        "HTTP 403",
		},
		{
			name:           "a 200 without windows is a real failure",
			responseStatus: http.StatusOK,
			responseBody:   `{"usage":{}}`,
			wantErr:        "no windows",
		},
		{
			name:           "an unparsable 200 is a real failure",
			responseStatus: http.StatusOK,
			responseBody:   `<html>maintenance</html>`,
			wantErr:        "no windows",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var gotPath, gotAuthorization, gotMethod string
			fakeAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotAuthorization, gotMethod = r.URL.Path, r.Header.Get("Authorization"), r.Method
				w.WriteHeader(testCase.responseStatus)
				_, _ = w.Write([]byte(testCase.responseBody))
			}))
			defer fakeAPIServer.Close()

			usageWindows, err := fetchUsageWindows(context.Background(), fakeAPIServer.Client(), fakeAPIServer.URL, "secret")

			if gotMethod != http.MethodGet || gotPath != "/usage" {
				t.Errorf("request = %s %s, want GET /usage", gotMethod, gotPath)
			}
			if gotAuthorization != "Bearer secret" {
				t.Errorf("Authorization = %q, want Bearer secret", gotAuthorization)
			}
			if testCase.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), testCase.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, testCase.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(usageWindows) != testCase.wantWindowCount {
				t.Errorf("usageWindows = %d entries, want %d", len(usageWindows), testCase.wantWindowCount)
			}
			if got := usageWindows["monthly"].Status; got != rateLimitedStatus {
				t.Errorf("monthly.Status = %q, want %q", got, rateLimitedStatus)
			}
		})
	}
}

func TestFetchUsageWindowsUnreachableAPI(t *testing.T) {
	closedServer := httptest.NewServer(http.NotFoundHandler())
	closedServer.Close()

	_, err := fetchUsageWindows(context.Background(), http.DefaultClient, closedServer.URL, "secret")
	if err == nil || !strings.Contains(err.Error(), "usage request") {
		t.Fatalf("err = %v, want a usage request error", err)
	}
}

func TestDescribeLimitedWindows(t *testing.T) {
	okWindow := UsageWindow{Status: "ok", UsedPercent: 40, ResetsAt: "2026-01-06T00:00:00.000Z"}
	limitedMonthlyWindow := UsageWindow{Status: rateLimitedStatus, UsedPercent: 100, ResetsAt: "2026-02-01T00:00:00.000Z"}
	limitedRollingWindow := UsageWindow{Status: rateLimitedStatus, UsedPercent: 100, ResetsAt: "2026-01-01T05:00:00.000Z"}

	testCases := []struct {
		name         string
		usageWindows map[string]UsageWindow
		want         []string
	}{
		{
			name:         "all windows ok means nothing is limited",
			usageWindows: map[string]UsageWindow{"rolling": okWindow, "weekly": okWindow, "monthly": okWindow},
			want:         nil,
		},
		{
			name:         "a limited window is named with its reset time",
			usageWindows: map[string]UsageWindow{"rolling": okWindow, "weekly": okWindow, "monthly": limitedMonthlyWindow},
			want:         []string{"monthly window at 100%, resets at 2026-02-01T00:00:00.000Z"},
		},
		{
			name:         "every limited window is named, in a stable order",
			usageWindows: map[string]UsageWindow{"rolling": limitedRollingWindow, "weekly": okWindow, "monthly": limitedMonthlyWindow},
			want: []string{
				"monthly window at 100%, resets at 2026-02-01T00:00:00.000Z",
				"rolling window at 100%, resets at 2026-01-01T05:00:00.000Z",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := describeLimitedWindows(testCase.usageWindows); !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("describeLimitedWindows() = %q, want %q", got, testCase.want)
			}
		})
	}
}
