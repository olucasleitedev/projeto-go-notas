//go:build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func gatewayURL() string {
	if u := os.Getenv("GATEWAY_URL"); u != "" {
		return u
	}
	return "http://localhost:8080"
}

func TestHealthEndpoints(t *testing.T) {
	endpoints := []string{"/health"}
	for _, path := range endpoints {
		res, err := http.Get(gatewayURL() + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		_ = res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET %s: status %d", path, res.StatusCode)
		}
	}
}

func TestCreateNote_AuditEvent(t *testing.T) {
	body, _ := json.Marshal(map[string]string{
		"title":   "Integration test",
		"content": "kafka + postgres",
	})

	res, err := http.Post(gatewayURL()+"/api/notes", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create note: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status: %d", res.StatusCode)
	}

	var created map[string]any
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created["id"] == nil || created["id"] == "" {
		t.Fatal("expected note id")
	}

	var events []map[string]any
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		auditRes, err := http.Get(gatewayURL() + "/api/audit/events")
		if err != nil {
			t.Fatalf("audit: %v", err)
		}
		if err := json.NewDecoder(auditRes.Body).Decode(&events); err != nil {
			auditRes.Body.Close()
			t.Fatal(err)
		}
		auditRes.Body.Close()

		for _, evt := range events {
			if evt["type"] == "note.created" && evt["note_id"] == created["id"] {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("audit event note.created not found for id %v", created["id"])
}
