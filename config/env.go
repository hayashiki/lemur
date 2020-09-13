package config

import (
	"fmt"
	"os"
)

var (
	keyPort          = "PORT"
	keyGCPProjectID  = "PROJECT_ID"
)

func NewReadMustFromEnv() (*Config, error) {
	cfg := &Config{}
	envs := getEnvs(keyPort, keyGCPProjectID)

	cfg.Port = envs[keyPort]
	if cfg.Port == "" {
		cfg.Port = "8000"
	}

	cfg.GCPProjectID = envs[keyGCPProjectID]
	if cfg.GCPProjectID == "" {
		return nil, fmt.Errorf("PROJECT_ID environment must be defined")
	}

	return cfg, nil
}

type Config struct {
	Port         string
	GCPProjectID string
}

func getEnvs(names ...string) map[string]string {
	envs := map[string]string{}
	for _, name := range names {
		envs[name] = os.Getenv(name)
	}
	return envs
}
