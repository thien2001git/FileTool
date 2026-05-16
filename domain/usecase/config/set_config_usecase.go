package config

import (
	"file_tool/domain/entity"
	"file_tool/domain/repository"
)

type SetConfigUseCase interface {
	SetConfigValue(key entity.ConfigKey, value string) error
}

type setConfigUseCase struct {
	configRepo repository.ConfigRepository
}

func NewSetConfigUseCase(configRepo repository.ConfigRepository) SetConfigUseCase {
	return &setConfigUseCase{
		configRepo: configRepo,
	}
}

func (uc *setConfigUseCase) SetConfigValue(key entity.ConfigKey, value string) error {
	return uc.configRepo.SetConfigValue(key, value)
}
