package httpCalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/imlogang/server-updater/internal/config"
)

type PteroErrorResponse struct {
	Errors []struct {
		Code   string `json:"code"`
		Status string `json:"status"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

type ClientServerResponse struct {
	Attributes struct {
		CurrentState *string `json:"current_state"`
	} `json:"attributes"`
}

func ReinstallServer(cfg *config.Config) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := fmt.Sprintf("%s/api/client/servers/%s/settings/reinstall", cfg.PteroURL, cfg.ServerID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.PteroToken))
	req.Header.Add("Accept", "Application/vnd.pterodactyl.v1+json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusAccepted {
		return fmt.Sprintf("The server with ID: %s was updated", cfg.ServerID), nil
	}

	// /settings/reisntall does not return a status code for a failed reinstall
	// should probably add some error handling if it does ever fail
	return "", nil
}

func PowerServer(cfg *config.Config, powerState string) (string, error) {
	body := map[string]string{
		"signal": powerState,
	}

	jsonStr, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/api/client/servers/%s/power", cfg.PteroURL, cfg.ServerID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.PteroToken))
	req.Header.Set("Accept", "application/vnd.pterodactyl.v1+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNoContent {
		return fmt.Sprintf("Power state %s initiated on %s", powerState, cfg.ServerID), nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("request failed: status=%d (failed reading body)", resp.StatusCode)
	}

	var apiErr PteroErrorResponse
	if err := json.Unmarshal(bodyBytes, &apiErr); err == nil && len(apiErr.Errors) > 0 {
		e := apiErr.Errors[0]
		return "", fmt.Errorf("%s: %s", e.Code, e.Detail)
	}

	return "", nil
}

func NotifyMinecraftServer(cfg *config.Config, command string) (string, error) {
	mcCommand := fmt.Sprintf("say %s", command)
	body := map[string]string{
		"command": mcCommand,
	}

	jsonStr, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/api/client/servers/%s/command", cfg.PteroURL, cfg.ServerID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.PteroToken))
	req.Header.Set("Accept", "application/vnd.pterodactyl.v1+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNoContent {
		return fmt.Sprintf("Command was issued on %s", cfg.ServerID), nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("request failed: status=%d (failed reading body)", resp.StatusCode)
	}

	var apiErr PteroErrorResponse
	if err := json.Unmarshal(bodyBytes, &apiErr); err == nil && len(apiErr.Errors) > 0 {
		e := apiErr.Errors[0]
		return "", fmt.Errorf("%s: %s", e.Code, e.Detail)
	}

	return "", nil
}

func GetStatus(cfg *config.Config) (*string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("%s/api/client/servers/%s/resources", cfg.PteroURL, cfg.ServerID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PteroToken)
	req.Header.Set("Accept", "application/vnd.pterodactyl.v1+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var parsed ClientServerResponse
	err = json.NewDecoder(resp.Body).Decode(&parsed)
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		return nil, err
	}
	return parsed.Attributes.CurrentState, nil
}

func WaitForOfflineThenStartMaybe(cfg *config.Config) (string, error) {
	for {
		status, err := GetStatus(cfg)
		if err != nil {
			return "", err
		}
		if status == nil && cfg.ServerID == "a6615eb7" {
			break
		}

		if *status == "running" {
			break
		}

		if *status == "installing" {
			time.Sleep(3 * time.Second)
			continue
		}

		time.Sleep(3 * time.Second)
	}
	// The Minecraft server(s) need to be 'turned on' after a reinstall, so we call PowerServer conditionally.
	// The Satisfactory server just needs to be restarted to be updated, but want to post a final message that it's online.
	if cfg.ServerID == "a6615eb7" {
		_, err := PowerServer(cfg, "start")
		if err != nil {
			return "", fmt.Errorf("there was an issue starting the %s server: %s", cfg.ServerName, err)
		}
	}
	
	resp := fmt.Sprintf("The %s server is back online.", cfg.ServerName)
	return resp, nil
}
