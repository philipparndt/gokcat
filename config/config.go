package config

import (
	"encoding/json"
	"github.com/philipparndt/go-logger"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	Broker         string `json:"broker"`
	SchemaRegistry struct {
		Url      string `json:"url"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Insecure bool   `json:"insecure,omitempty"`
	} `json:"schemaRegistry"`
	Certs struct {
		Ca         string `json:"ca"`
		ClientCert string `json:"clientCert"`
		ClientKey  string `json:"clientKey"`
		Insecure   bool   `json:"insecure,omitempty"`
	} `json:"certs"`
	LogLevel string `json:"logLevel"`
}

func ReplaceEnvVariables(input []byte) []byte {
	envVariableRegex := regexp.MustCompile(`\${([^}]+)}`)

	return envVariableRegex.ReplaceAllFunc(input, func(match []byte) []byte {
		envVarName := match[2 : len(match)-1] // Extract the variable name without "${}".
		return []byte(os.Getenv(string(envVarName)))
	})
}

func LoadConfig(file string) (Config, error) {
	configFile := expandPath(file)
	data, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error("Error reading config file", err)
		return Config{}, err
	}

	data = ReplaceEnvVariables(data)

	// Create a Config object
	var cfg Config

	// Unmarshal the JSON data into the Config object
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("Unmarshalling JSON:", err)
		return Config{}, err
	}

	updateCertsPath(configFile, &cfg)

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg, nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return strings.Replace(path, "~", home, 1)
		}
	}
	return path
}

func updateCertsPath(file string, cfg *Config) {
	configDir := ""
	absFile, err := filepath.Abs(file)
	if err == nil {
		configDir = filepath.Dir(absFile)
	} else {
		configDir = filepath.Dir(file)
	}
	if cfg.Certs.Ca != "" && !filepath.IsAbs(cfg.Certs.Ca) {
		cfg.Certs.Ca = filepath.Join(configDir, cfg.Certs.Ca)
	}
	if cfg.Certs.ClientCert != "" && !filepath.IsAbs(cfg.Certs.ClientCert) {
		cfg.Certs.ClientCert = filepath.Join(configDir, cfg.Certs.ClientCert)
	}
	if cfg.Certs.ClientKey != "" && !filepath.IsAbs(cfg.Certs.ClientKey) {
		cfg.Certs.ClientKey = filepath.Join(configDir, cfg.Certs.ClientKey)
	}
}
