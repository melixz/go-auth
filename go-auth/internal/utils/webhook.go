package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type WebhookPayload struct {
	UserID    string `json:"user_id"`
	OldIP     string `json:"old_ip"`
	NewIP     string `json:"new_ip"`
	UserAgent string `json:"user_agent"`
	Timestamp int64  `json:"timestamp"`
}

func SendWebhook(userID, oldIP, newIP, userAgent string) {
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return
	}
	payload := WebhookPayload{
		UserID:    userID,
		OldIP:     oldIP,
		NewIP:     newIP,
		UserAgent: userAgent,
		Timestamp: time.Now().Unix(),
	}
	jsonData, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
}
