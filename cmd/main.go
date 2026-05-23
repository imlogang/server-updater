package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/imlogang/server-updater/internal/config"
	"github.com/imlogang/server-updater/internal/servers"
)

func lookupWebhook(key string) string {
	webhook, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s is not set.", key)
	}
	return webhook
}

func main() {
	cfg := config.Config{}
	kong.Parse(&cfg)
	switch cfg.ServerID {
	// SMP Vanilla
	case "a6615eb7":
		cfg.DiscordWebhookLink = lookupWebhook("DISCORD_WEBHOOK_LINK_MINECRAFT")
		cfg.ServerName = "SMP Vanilla"

		err := servers.SMPVanilla(&cfg)
		if err != nil {
			panic(err)
		}
	//satisfactory
	case "6b774df5":
		cfg.DiscordWebhookLink = lookupWebhook("DISCORD_WEBHOOK_LINK_SATISFACTORY")
		cfg.ServerName = "Satisfactory"

		err := servers.Satisfactory(&cfg)
		if err != nil {
			panic(err)
		}
	}
}
