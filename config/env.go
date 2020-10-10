package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port          string `envconfig:"PORT" default:"8080"`
	GCPProjectID  string `envconfig:"PROJECT_ID" require:"true"`
	GCPLocationID string `envconfig:"LOCATION_ID" require:"true"`
	DocBaseTeam   string `envconfig:"DOCBASE_TEAM" require:"true"`
	DocBaseToken  string `envconfig:"DOCBASE_TOKEN" require:"true"`
	GithubSecret  string `envconfig:"GITHUB_SECRET" require:"true"`
	GithubOrg     string `envconfig:"GITHUB_ORG" require:"true"`
	GithubRepo    string `envconfig:"GITHUB_REPO" require:"true"`
}

func NewReadMustFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
