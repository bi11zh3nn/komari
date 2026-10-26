package bark

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/komari-monitor/komari/utils/messageSender"
)

func TestBarkSenderUsesLegacyConfiguration(t *testing.T) {
	var payload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/push" {
			t.Fatalf("path = %q, want /push", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"code": 200})
	}))
	defer server.Close()

	err := (&BarkSender{}).Send(context.Background(), messageSender.Notification{
		Title:   "Offline",
		Message: "VMISS TRI is offline",
	}, map[string]any{
		"server_url": server.URL,
		"device_key": "device",
		"icon":       "https://example.com/icon.png",
		"level":      "timeSensitive",
	})
	if err != nil {
		t.Fatal(err)
	}
	if payload["device_key"] != "device" || payload["title"] != "Offline" || payload["body"] != "VMISS TRI is offline" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}
