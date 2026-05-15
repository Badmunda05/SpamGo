package client

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"pbxgo/config"
	"pbxgo/database"
	"pbxgo/modules"
	"syscall"

	"github.com/amarnathcjd/gogram/telegram"
)

var Clients []*telegram.Client

func InitMultiBots() error {
	tokens := config.GetBotTokens()
	if len(tokens) == 0 {
		return fmt.Errorf("no BOT_TOKEN found in .env")
	}

	for i, token := range tokens {
		cfg := telegram.ClientConfig{
			AppID:       config.AppConfig.AppID,
			AppHash:     config.AppConfig.AppHash,
			LogLevel:    telegram.LogInfo,
			SessionName: fmt.Sprintf("bot_%d", i+1),
			DeviceConfig: telegram.DeviceConfig{
				DeviceModel:    fmt.Sprintf("PbxGo-Bot%d", i+1),
				SystemVersion:  "Go 1.26",
				AppVersion:     "v2.1.0",
				LangCode:       "en",
				SystemLangCode: "en-US",
			},
		}

		cl, err := telegram.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("bot %d creation failed: %w", i+1, err)
		}

		if _, err = cl.Conn(); err != nil {
			slog.Error("Bot connect failed", "bot", i+1, "error", err)
			continue
		}

		// If session already exists, reuse it — do not log out
		// Logging out on restart was causing the bot to hang on reconnect
		if authorized, _ := cl.IsAuthorized(); !authorized {
			if err = cl.LoginBot(token); err != nil {
				slog.Error("Bot login failed", "bot", i+1, "error", err)
				continue
			}
		} else {
			slog.Info("Existing session restored", "bot", i+1)
		}

		Clients = append(Clients, cl)
		slog.Info("Bot started", "bot", i+1, "id", cl.Me().ID)

		go modules.ResumeRaids(cl)
	}

	return nil
}

func RegisterHandlers() {
	for _, cl := range Clients {
		modules.Load(cl)
	}
	slog.Info("Handlers registered", "bots", len(Clients))
}

func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		slog.Info("Shutting down PbxGo...")
		database.Disconnect()
		os.Exit(0)
	}()

	for _, cl := range Clients {
		go cl.Idle()
	}

	select {}
}
