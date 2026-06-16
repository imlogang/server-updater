package servers

import (
	"fmt"

	"github.com/imlogang/server-updater/internal/config"
	"github.com/imlogang/server-updater/internal/httpCalls"
)

func UpdateDiscord(cfg *config.Config, time string, clock string, latestVersion String) error {
	message := fmt.Sprintf("The server will be updated in %s %s to %s. Please update your client!", time, clock)
	resp, err := httpCalls.NotifyDiscord(cfg, message)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}

func UpdateServer(cfg *config.Config, time string, clock string) error {
	message := fmt.Sprintf("The server will be updated in %s %s. Please update your client!", time, clock)
	resp, err := httpCalls.NotifyMinecraftServer(cfg, message)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}

func UpdateDiscordFinal(cfg *config.Config) (error) {
	message := fmt.Sprintf("The %s server has been updated.", cfg.ServerName)
	resp, err := httpCalls.NotifyDiscord(cfg, message)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}