package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type SiteConfig struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type S3Config struct {
	Endpoint        string `yaml:"endpoint"`
	Region          string `yaml:"region"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

type Config struct {
	Port int        `yaml:"port"`
	Site SiteConfig `yaml:"site"`
	S3   S3Config   `yaml:"s3"`
}

// Load reads configuration from a YAML file and applies environment variable overrides.
// Environment variables take precedence over file values.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Port: 3000, // default port
		Site: SiteConfig{
			Title: "S3 Bucket Browser", // default title
		},
	}

	// Try to load from file if it exists
	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Apply environment variable overrides
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if endpoint := os.Getenv("S3_ENDPOINT"); endpoint != "" {
		cfg.S3.Endpoint = endpoint
	}

	if region := os.Getenv("S3_REGION"); region != "" {
		cfg.S3.Region = region
	}

	if bucket := os.Getenv("S3_BUCKET"); bucket != "" {
		cfg.S3.Bucket = bucket
	}

	if accessKey := os.Getenv("S3_ACCESS_KEY_ID"); accessKey != "" {
		cfg.S3.AccessKeyID = accessKey
	}

	if secretKey := os.Getenv("S3_SECRET_ACCESS_KEY"); secretKey != "" {
		cfg.S3.SecretAccessKey = secretKey
	}

	if title := os.Getenv("SITE_TITLE"); title != "" {
		cfg.Site.Title = title
	}

	if description := os.Getenv("SITE_DESCRIPTION"); description != "" {
		cfg.Site.Description = description
	}

	return cfg, nil
}
