package config

import (
	"file_tool/domain/entity"
	"file_tool/domain/repository"
)

type GetConfigUseCase interface {
	GetConfigValue(key entity.ConfigKey) (string, error)
}

type getConfigUseCase struct {
	configRepo repository.ConfigRepository
}

func NewGetConfigUseCase(configRepo repository.ConfigRepository) GetConfigUseCase {
	return &getConfigUseCase{
		configRepo: configRepo,
	}
}

func (uc *getConfigUseCase) GetConfigValue(key entity.ConfigKey) (string, error) {
	return uc.configRepo.GetConfigValue(key)
}
