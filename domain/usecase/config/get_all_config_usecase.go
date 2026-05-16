package config

import (
	"file_tool/domain/entity"
	"file_tool/domain/repository"
)

type GetAllConfigUseCase interface {
	GetAllConfigValues() (map[entity.ConfigKey]string, error)
}

type getAllConfigUseCase struct {
	configRepo repository.ConfigRepository
}

func NewGetAllConfigUseCase(configRepo repository.ConfigRepository) GetAllConfigUseCase {
	return &getAllConfigUseCase{
		configRepo: configRepo,
	}
}

func (uc *getAllConfigUseCase) GetAllConfigValues() (map[entity.ConfigKey]string, error) {
	return uc.configRepo.GetAllConfigValues()
}
