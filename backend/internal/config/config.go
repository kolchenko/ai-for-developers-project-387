package config

import "os"

type Config struct {
	Addr   string
	DBPath string
	WebDir string
}

func Load() Config {
	return Config{
		Addr:   addr(),
		DBPath: envOr("DB_PATH", "data.db"),
		WebDir: envOr("WEB_DIR", ""),
	}
}

func addr() string {
	if v := os.Getenv("ADDR"); v != "" {
		return v
	}
	return "0.0.0.0:" + envOr("PORT", "4010")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
