package repository

import (
	"file_tool/domain/entity"
)


type ConfigRepository interface {
	// GetConfigValue retrieves a configuration value based on the provided key.
	GetConfigValue(key entity.ConfigKey) (string, error)
	SetConfigValue(key entity.ConfigKey, value string) (error)
	GetAllConfigValues() (map[entity.ConfigKey]string, error)
}