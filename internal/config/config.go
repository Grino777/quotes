package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type SQLiteConfig struct {
	Addr string `yaml:"local_path" default:"storage/quotes.sqlite"`
}

type APIConfig struct {
	Addr string `yaml:"addr" env-default:"127.0.0.1"`
	Port string `yaml:"port" default:"8090"`
}

type Config struct {
	SQLite  SQLiteConfig `yaml:"sqlite" required:"true"`
	API     APIConfig    `yaml:"api" required:"true"`
	BaseDir string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := getBaseDir(cfg); err != nil {
		return cfg, fmt.Errorf("failed to get base dir")
	}

	configFile := filepath.Join(cfg.BaseDir, "configs/local.yml")
	if err := cleanenv.ReadConfig(configFile, cfg); err != nil {
		return cfg, err
	}

	dbPath := filepath.Join(cfg.BaseDir, cfg.SQLite.Addr)
	cfg.SQLite.Addr = dbPath

	return cfg, nil
}

func getBaseDir(cfg *Config) error {
	// Получаем путь к исполняемому файлу
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	baseDir := filepath.Dir(exePath)
	baseDir = filepath.Join(baseDir, "../..")

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cfg.BaseDir = absBaseDir
	return nil
}
