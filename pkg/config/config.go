package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	WorkspaceRoot   string          `yaml:"workspace_root"`
	PeriodicNotes   PeriodicNotes   `yaml:"periodic_notes"`
	WeeklyView      WeeklyViewConfig `yaml:"weekly_view"`
}

type PeriodicNotes struct {
	WeeklySubdir     string `yaml:"weekly_subdir"`
	WeeklyNameFormat string `yaml:"weekly_name_format"`
}

type WeeklyViewConfig struct {
	EnabledComponents WeeklyViewComponents `yaml:"enabled_components"`
}

type WeeklyViewComponents struct {
	HabitTracker bool `yaml:"habit_tracker"`
}

func NewConfigWithDefaults() *Config {
	return &Config{
		PeriodicNotes: PeriodicNotes{
			WeeklySubdir:     "_periodic/weekly",
			WeeklyNameFormat: "YYYY-[W]WW",
		},
		WeeklyView: WeeklyViewConfig{
			EnabledComponents: WeeklyViewComponents{
				HabitTracker: true,
			},
		},
	}
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".notedown", "planner")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

func ValidateWorkspacePath(path string) error {
	if path == "" {
		return os.ErrNotExist
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Check if it's a directory
	if !info.IsDir() {
		return os.ErrInvalid
	}

	// Check if we can write to it
	testFile := filepath.Join(path, ".notedown_test")
	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	file.Close()
	os.Remove(testFile)

	return nil
}
