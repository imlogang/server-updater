package servers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/imlogang/server-updater/internal/config"
	"github.com/imlogang/server-updater/internal/httpCalls"
)

type MojangResponse struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
}

func getMinecraftVersion() (string, error) {
	url := "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var manifest MojangResponse
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return "", err
	}

	return manifest.Latest.Release, nil

}

func SMPVanilla(cfg *config.Config) error {
	println("Running commands for SMP Vanilla")
	latestVersion, err := getMinecraftVersion()
	if err != nil {
		return fmt.Errorf("there was an error fetching the latest Minecraft version: %s", err)
	}
	cfg.LatestVersion = latestVersion

	err = updateServerandDiscord(cfg)
	if err != nil {
		return err
	}

	resp, err := httpCalls.ReinstallServer(cfg)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	resp, err = httpCalls.WaitForOfflineThenStartMaybe(cfg)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	err = UpdateDiscordFinal(cfg)
	if err != nil {
		return err
	}

	return nil
}

func updateServerandDiscord(cfg *config.Config) error {
	err := UpdateDiscord(cfg, "5", "minutes")
	if err != nil {
		return err
	}
	err = UpdateServer(cfg, "5", "minutes")
	if err != nil {
		return err
	}
	time.Sleep(4 * time.Minute)

	err = UpdateDiscord(cfg, "1", "minute")
	if err != nil {
		return err
	}
	err = UpdateServer(cfg, "1", "minute")
	if err != nil {
		return err
	}
	time.Sleep(55 * time.Second)

	err = UpdateDiscord(cfg, "5", "seconds")
	if err != nil {
		return err
	}
	err = UpdateServer(cfg, "1", "seconds")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	return nil
}
