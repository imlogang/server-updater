package servers

import (
	"fmt"

	"github.com/imlogang/server-updater/internal/config"
	"github.com/imlogang/server-updater/internal/httpCalls"
)

func UpdateDiscord(cfg *config.Config, time string, clock string) error {
	message := fmt.Sprintf("The server will be updated in %s clock to %s. Please update your client!", time, clock)
	resp, err := httpCalls.NotifyDiscord(cfg.DiscordWebhookLink, message)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}

func UpdateServer(cfg *config.Config, time string, clock string) error {
	message := fmt.Sprintf("The server will be updated in %s seconds to %s. Please update your client!", time, clock)
	resp, err := httpCalls.NotifyDiscord(cfg.DiscordWebhookLink, message)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}
