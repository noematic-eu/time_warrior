// Package configuration provides some default directory and file names.
package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

// Config for the application files/folders
type Config struct {
	homeDirectory   string
	dataFolder      string
	pendingFilename string
	projectFile     string
	feeFile         string
}

// New returns a new configuration with some sane defaults
func New() *Config {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &Config{
		homeDirectory:   home,
		dataFolder:      "time_warrior",
		pendingFilename: ".pending",
		projectFile:     ".project",
		feeFile:         ".fees",
	}
}

func (c Config) DataDirectoryPath() string {
	return path.Join(c.homeDirectory, c.dataFolder)
}

func (c Config) PendingFilePath() string {
	return path.Join(c.DataDirectoryPath(), c.pendingFilename)
}

func (c Config) ProjectFilePath() string {
	return path.Join(c.DataDirectoryPath(), c.projectFile)
}

func (c Config) FeeFilePath() string {
	return path.Join(c.DataDirectoryPath(), c.feeFile)
}

func (c Config) VerifyDataFilesPresent() bool {
	if _, err := os.Stat(c.DataDirectoryPath()); err != nil {
		return false
	}

	if _, err := os.Stat(c.PendingFilePath()); err != nil {
		return false
	}

	if _, err := os.Stat(c.ProjectFilePath()); err != nil {
		return false
	}

	if _, err := os.Stat(c.FeeFilePath()); err != nil {
		return false
	}

	return true
}

// GetCurrentProject returns the current project or empty string if not set
func (c Config) GetCurrentProject() (string, error) {
	if _, err := os.Stat(c.ProjectFilePath()); err != nil {
		return "", nil
	}

	data, err := os.ReadFile(c.ProjectFilePath())
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// SetCurrentProject sets the current project
func (c Config) SetCurrentProject(project string) error {
	return os.WriteFile(c.ProjectFilePath(), []byte(project), 0644)
}

// GetProjectFee returns the fee per hour for a project
func (c Config) GetProjectFee(project string) (float64, error) {
	feePath := path.Join(c.DataDirectoryPath(), c.feeFile)
	if _, err := os.Stat(feePath); err != nil {
		return 0, nil
	}

	data, err := os.ReadFile(feePath)
	if err != nil {
		return 0, err
	}

	var fees map[string]float64
	if err := json.Unmarshal(data, &fees); err != nil {
		return 0, err
	}

	if fee, ok := fees[project]; ok {
		return fee, nil
	}
	return 0, nil
}

// SetProjectFee sets the fee per hour for a project
func (c Config) SetProjectFee(project string, fee float64) error {
	feePath := path.Join(c.DataDirectoryPath(), c.feeFile)

	// Read existing fees
	var fees map[string]float64
	if _, err := os.Stat(feePath); err == nil {
		data, err := os.ReadFile(feePath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &fees); err != nil {
			return err
		}
	} else {
		fees = make(map[string]float64)
	}

	// Update fee
	fees[project] = fee

	// Write back to file
	data, err := json.Marshal(fees)
	if err != nil {
		return err
	}

	return os.WriteFile(feePath, data, 0644)
}
