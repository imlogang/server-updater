package config

type Config struct {
	PteroToken         string `env:"PTERO_TOKEN" help:"API Key for Pterodactyl Instance"`
	PteroURL           string `env:"PTERO_URL" help:"Pterodactyl Instance URL"`
	ServerID           string `env:"SERVER_ID" help:"ID of the Pterodactyl Server"`
	DiscordWebhookLink string
	LatestVersion      string
	ServerName         string
}
