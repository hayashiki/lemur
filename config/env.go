package config

import (
	"fmt"
	"os"
)

var (
	keyPort         = "PORT"
	keyGCPProjectID = "PROJECT_ID"
	keyGCPLocationID = "LOCATION_ID"
	keyDocBaseTeam = "DOCBASE_TEAM"
	keyDocBaseToken = "DOCBASE_TOKEN"
)

func NewReadMustFromEnv() (*Config, error) {
	cfg := &Config{}
	envs := getEnvs(keyPort, keyGCPProjectID, keyGCPLocationID, keyDocBaseTeam, keyDocBaseToken)

	cfg.Port = envs[keyPort]
	if cfg.Port == "" {
		cfg.Port = "8000"
	}

	cfg.GCPProjectID = envs[keyGCPProjectID]
	if cfg.GCPProjectID == "" {
		return nil, fmt.Errorf("PROJECT_ID environment must be defined")
	}

	cfg.GCPLocationID = envs[keyGCPLocationID]
	if cfg.GCPLocationID == "" {
		return nil, fmt.Errorf("LOCATION_ID environment must be defined")
	}

	cfg.DocBaseTeam = envs[keyDocBaseTeam]
	if cfg.DocBaseTeam == "" {
		return nil, fmt.Errorf("DOCBASE_TEAM environment must be defined")
	}

	cfg.DocBaseToken = envs[keyDocBaseToken]
	if cfg.DocBaseToken == "" {
		return nil, fmt.Errorf("DOCBASE_TOKEN environment must be defined")
	}

	return cfg, nil
}

type Config struct {
	Port         string
	GCPProjectID string
	GCPLocationID string
	DocBaseTeam  string
	DocBaseToken string
}

func getEnvs(names ...string) map[string]string {
	envs := map[string]string{}
	for _, name := range names {
		envs[name] = os.Getenv(name)
	}
	return envs
}
