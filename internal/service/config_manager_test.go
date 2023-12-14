package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dough-calculator/internal/config"
)

func TestNewConfigManager(t *testing.T) {
	manager := NewConfigManager()

	assert.NotNil(t, manager)
}

func TestConfigManager_ParseConfig(t *testing.T) {
	t.Setenv("CONFIG_PATH", "testdata/application-config.yaml")

	manager := NewConfigManager()

	require.NoError(t, manager.ParseConfig())

	managerStr := manager.(*configManager)

	require.Equal(t, config.Config{
		Application: config.Application{
			Name: "dough-calculator",
			Rest: config.Rest{
				Server:               ":8080",
				ContextPath:          "/v1",
				ReadTimeout:          5,
				WriteTimeout:         5,
				IdleTimeout:          120,
				GraceShutdownTimeout: 20,
			},
		},
		Database: config.Database{
			Uri:               "mongodb://localhost:27017",
			ConnectionTimeout: 10000,
		},
	}, managerStr.config)
}

func TestConfigManager_ParseConfig_WithEnvironmentVariableOverride(t *testing.T) {
	t.Setenv("CONFIG_PATH", "testdata/application-config.yaml")
	t.Setenv("DATABASE_URI", "test_uri_override")

	manager := NewConfigManager()

	require.NoError(t, manager.ParseConfig())

	managerStr := manager.(*configManager)

	require.Equal(t, config.Config{
		Application: config.Application{
			Name: "dough-calculator",
			Rest: config.Rest{
				Server:               ":8080",
				ContextPath:          "/v1",
				ReadTimeout:          5,
				WriteTimeout:         5,
				IdleTimeout:          120,
				GraceShutdownTimeout: 20,
			},
		},
		Database: config.Database{
			Uri:               "test_uri_override",
			ConnectionTimeout: 10000,
		},
	}, managerStr.config)
}

func TestConfigManager_ParseConfig_InvalidConfigPath(t *testing.T) {
	t.Setenv("CONFIG_PATH", "testdata/invalid-config.yaml")

	manager := NewConfigManager()

	err := manager.ParseConfig()

	require.Error(t, err)
}

func TestConfigManager_ParseConfig_InvalidContentType(t *testing.T) {
	dir := t.TempDir()
	file, err := os.Create(dir + "/invalid-config-content.yaml")
	require.NoError(t, err)
	err = os.WriteFile(file.Name(), []byte("hello"), 0644)
	require.NoError(t, err)

	t.Setenv("CONFIG_PATH", file.Name())

	manager := NewConfigManager()

	err = manager.ParseConfig()

	require.Error(t, err)
}

func TestConfigManager_ParseConfig_WithEmptyContent(t *testing.T) {
	dir := t.TempDir()
	file, err := os.Create(dir + "/empty-config-content.yaml")
	require.NoError(t, err)
	err = os.WriteFile(file.Name(), []byte(""), 0644)
	require.NoError(t, err)

	t.Setenv("CONFIG_PATH", file.Name())

	manager := NewConfigManager()

	err = manager.ParseConfig()

	require.ErrorContains(t, err, fmt.Sprintf("failed to parse config file '%s'", file.Name()))
}

func TestConfigManager_GetConfig(t *testing.T) {
	t.Setenv("CONFIG_PATH", "testdata/application-config.yaml")

	manager := NewConfigManager()

	require.NoError(t, manager.ParseConfig())

	require.Equal(t, config.Config{
		Application: config.Application{
			Name: "dough-calculator",
			Rest: config.Rest{
				Server:               ":8080",
				ContextPath:          "/v1",
				ReadTimeout:          5,
				WriteTimeout:         5,
				IdleTimeout:          120,
				GraceShutdownTimeout: 20,
			},
		},
		Database: config.Database{
			Uri:               "mongodb://localhost:27017",
			ConnectionTimeout: 10000,
		},
	}, manager.GetConfig())
}
