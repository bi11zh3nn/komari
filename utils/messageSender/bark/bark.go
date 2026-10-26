package bark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/komari-monitor/komari/utils/messageSender"
)

const ChannelID = "bark"

type BarkSender struct{}

func (b *BarkSender) Send(ctx context.Context, notification messageSender.Notification, config map[string]any) error {
	var addition Addition
	if err := messageSender.DecodeConfiguration(config, &addition); err != nil {
		return err
	}
	if strings.TrimSpace(addition.DeviceKey) == "" {
		return fmt.Errorf("device key is required")
	}
	if strings.TrimSpace(notification.Message) == "" {
		return fmt.Errorf("message is empty")
	}
	serverURL := strings.TrimRight(strings.TrimSpace(addition.ServerURL), "/")
	if serverURL == "" {
		serverURL = "https://api.day.app"
	}
	parsed, err := url.Parse(serverURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("invalid server URL %q", serverURL)
	}

	payload := map[string]any{
		"body":       notification.Message,
		"device_key": addition.DeviceKey,
	}
	if notification.Title != "" {
		payload["title"] = notification.Title
	}
	if addition.Icon != "" {
		payload["icon"] = addition.Icon
	}
	if addition.Level != "" {
		payload["level"] = addition.Level
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal Bark payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/push", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create Bark request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("send Bark request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bark API returned status %d", resp.StatusCode)
	}
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result.Code != 0 && result.Code != 200 {
		return fmt.Errorf("Bark API error (code %d): %s", result.Code, result.Message)
	}
	return nil
}

func (*BarkSender) Unload() error { return nil }

func Register() error {
	return messageSender.RegisterNotificationChannel(
		ChannelID,
		messageSender.ManagedConfiguration("Bark", Addition{}),
		&BarkSender{},
	)
}

var _ messageSender.NotificationChannel = (*BarkSender)(nil)
