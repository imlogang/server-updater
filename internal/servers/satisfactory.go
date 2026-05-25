package servers

import (
	"fmt"
	"time"

	"github.com/imlogang/server-updater/internal/config"
	"github.com/imlogang/server-updater/internal/httpCalls"
)

func Satisfactory(cfg *config.Config) error {
	println("Running commands for Satisfactory")
	err := updateServer(cfg)
	if err != nil {
		return err
	}
	resp, err := httpCalls.PowerServer(cfg, "restart")
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

func updateServer(cfg *config.Config) error {
	err := UpdateDiscord(cfg, "5", "minutes")
	if err != nil {
		return err
	}
	time.Sleep(4 * time.Minute)

	err = UpdateDiscord(cfg, "1", "minute")
	if err != nil {
		return err
	}
	time.Sleep(55 * time.Second)

	err = UpdateDiscord(cfg, "1", "seconds")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	return nil
}
