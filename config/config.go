package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppID    int32
	AppHash  string
	OwnerID  int64
	MongoURL string

	// Images
	StartPic string
	HelpPic  string
}

var AppConfig Config

func Load() {

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file, using system environment")
	}

	appID, _ := strconv.ParseInt(
		mustGetEnv("APP_ID"),
		10,
		32,
	)

	ownerID, _ := strconv.ParseInt(
		mustGetEnv("OWNER_ID"),
		10,
		64,
	)

	AppConfig = Config{
		AppID:    int32(appID),
		AppHash:  mustGetEnv("APP_HASH"),
		OwnerID:  ownerID,
		MongoURL: getEnv("MONGO_URL", ""),

		// Images
		StartPic: getEnv(
			"START_PIC",
			"https://files.tgvibes.online/5JreGgKB.jpg",
		),

		HelpPic: getEnv(
			"HELP_PIC",
			"https://files.tgvibes.online/5JreGgKB.jpg",
		),
	}

	slog.Info(
		"✅ Config Loaded",
		"bots",
		len(GetBotTokens()),
		"owner",
		ownerID,
	)
}

func GetBotTokens() []string {

	var tokens []string

	for i := 1; i <= 50; i++ {

		if t := getEnv(
			"BOT_TOKEN"+strconv.Itoa(i),
			"",
		); t != "" {

			tokens = append(tokens, t)
		}
	}

	return tokens
}

func getEnv(key, fallback string) string {

	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func mustGetEnv(key string) string {

	v := os.Getenv(key)

	if v == "" {

		slog.Error(
			"Missing required env",
			"key",
			key,
		)

		os.Exit(1)
	}

	return v
}
