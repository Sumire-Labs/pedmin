package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/disgoorg/snowflake/v2"
)

type Config struct {
	Token            string
	AppID            snowflake.ID
	LavalinkHost     string
	LavalinkPassword string
	DataDir          string
	DBPath           string
}

func Load() (*Config, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}

	appIDStr := os.Getenv("DISCORD_APP_ID")
	if appIDStr == "" {
		return nil, fmt.Errorf("DISCORD_APP_ID is required")
	}
	appID, err := snowflake.Parse(appIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DISCORD_APP_ID: %w", err)
	}

	lavalinkHost := os.Getenv("LAVALINK_HOST")
	if lavalinkHost == "" {
		lavalinkHost = "lavalink:2333"
	}

	lavalinkPassword := os.Getenv("LAVALINK_PASSWORD")
	if lavalinkPassword == "" {
		lavalinkPassword = "youshallnotpass"
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(dataDir, "pedmin.db")
	}

	return &Config{
		Token:            token,
		AppID:            appID,
		LavalinkHost:     lavalinkHost,
		LavalinkPassword: lavalinkPassword,
		DataDir:          dataDir,
		DBPath:           dbPath,
	}, nil
}
