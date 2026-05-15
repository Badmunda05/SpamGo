package main

import (
	"log/slog"
	"os"
	"pbxgo/client"
	"pbxgo/config"
	"pbxgo/database"

	_ "pbxgo/modules"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	slog.Info("🚀 PbxGo-Multi Starting...", "version", "v2.1.0")

	config.Load()
	database.Init()
	database.LoadSudoUsers()

	if err := client.InitMultiBots(); err != nil {
		slog.Error("Failed to start bots", "error", err)
		os.Exit(1)
	}

	client.RegisterHandlers()

	slog.Info("✅ All bots are running!")
	client.Run()
}
