package main

import (
	"awesome-piracy-bot/pkg/discord"
	"awesome-piracy-bot/pkg/telegram"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

var goroutineDelta = make(chan int)

type Config struct {
	Telegram struct {
		Enabled  bool
		APIToken string
	}
	Discord struct {
		Enabled  bool
		APIToken string
	}
	Reddit struct {
		Enabled bool
	}
	IRC struct {
		Enabled bool
	}
}

type WatcherConfig struct {
	Type     string
	Enabled  bool
	APIToken string
}

type WatcherConfigs []WatcherConfig

func main() {
	// get configuration
	var name = "awesome-piracy-bot"
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", name))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.config/%s/", name))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Fatal error in config file: %s \n", err)
	}

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Panicf("Error parsing config file, %v \n", err)
	}

	discordConfig := WatcherConfig{
		Type:     "Discord",
		Enabled:  config.Discord.Enabled,
		APIToken: config.Discord.APIToken,
	}

	telegramConfig := WatcherConfig{
		Type:     "Telegram",
		Enabled:  config.Telegram.Enabled,
		APIToken: config.Telegram.APIToken,
	}

	watcherConfigs := WatcherConfigs{
		telegramConfig,
		discordConfig,
	}

	// start watchers
	for _, w := range watcherConfigs {
		go w.startWatcher()
	}

	// TODO: send URLs back to main() via a channel
	// TODO: add elasticsearch URL destination with metadata
	// TODO: add metadata to URLs, e.g. HTTP response, HTML title, protocol

	numGoroutines := 0
	for diff := range goroutineDelta {
		numGoroutines += diff
		if numGoroutines == 0 {
			os.Exit(0)
		}
	}
}

func (c WatcherConfig) startWatcher() {
	if c.Enabled != true {
		log.Printf("[INFO] %s disabled - skipping", c.Type)
	} else {
		if c.Type == "Telegram" {
			telegram.Run(c.APIToken)
		}
		if c.Type == "Discord" {
			discord.Run(c.APIToken)
		}
	}
	goroutineDelta <- +1
	go f()
}

func f() {
	goroutineDelta <- -1
}
