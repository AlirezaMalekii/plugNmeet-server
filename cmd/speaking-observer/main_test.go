package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSendJournalSignsExactBodyAndKeepsStableID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session-segment.jsonl")
	if err := os.WriteFile(path, []byte("{\"event\":\"speech_started\",\"room_id\":\"room\",\"room_sid\":\"sid\",\"participant_key\":\"user\",\"timestamp_ms\":1000}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var ids []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		ids = append(ids, body["id"].(string))
		parts := strings.Split(r.Header.Get("Authorization"), ".")
		if len(parts) != 3 {
			t.Error("missing JWT")
			w.WriteHeader(401)
			return
		}
		claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			t.Error(err)
		}
		var claims map[string]any
		if err = json.Unmarshal(claimsBytes, &claims); err != nil {
			t.Error(err)
		}
		if claims["iss"] != "key" || claims["sha256"] == "" {
			t.Error("missing issuer or body checksum")
		}
		if body["event"] != "ineshast_speaking_events" {
			t.Error("wrong event")
		}
		w.WriteHeader(202)
	}))
	defer server.Close()
	for i := 0; i < 2; i++ {
		if err := sendJournal(context.Background(), server.Client(), path, server.URL, "key", "secret"); err != nil {
			t.Fatal(err)
		}
	}
	if len(ids) != 2 || ids[0] != ids[1] {
		t.Fatal("retry must preserve deduplication ID")
	}
}
func TestSendJournalRejectsUnacknowledgedDelivery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.jsonl")
	os.WriteFile(path, []byte("{\"room_id\":\"room\",\"room_sid\":\"sid\"}\n"), 0600)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(503) }))
	defer server.Close()
	if err := sendJournal(context.Background(), server.Client(), path, server.URL, "key", "secret"); err == nil {
		t.Fatal("503 must remain retryable")
	}
}
func TestRecoverJournalClosesAtLastObservedTime(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	os.WriteFile(path, []byte("{\"event\":\"speech_started\",\"room_id\":\"room\",\"room_sid\":\"sid\",\"timestamp_ms\":1000}\n"), 0600)
	if err := recoverJournals(dir); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	var last map[string]any
	if err = json.Unmarshal([]byte(lines[len(lines)-1]), &last); err != nil {
		t.Fatal(err)
	}
	if last["timestamp_ms"] != float64(1000) || last["observation_gap"] != true {
		t.Fatal("recovery must identify unobserved gap without inflating speech")
	}
	if _, err = os.Stat(path + ".ready"); err != nil {
		t.Fatal(err)
	}
	if err = recoverJournals(dir); err != nil {
		t.Fatal(err)
	}
	again, _ := os.ReadFile(path)
	if string(again) != string(body) {
		t.Fatal("recovery must be idempotent")
	}
}
