package httpCalls

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func NotifyDiscord(discordLink string, message string) (string, error) {
	body := map[string]string{
		"content": message,
	}

	jsonStr, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("POST", discordLink, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNoContent {
		return "Notified Discord", nil
	}
	return "", nil
}
