package provider

import (
	"file_tool/data/repository_impl"
	"file_tool/domain/repository"
	"file_tool/domain/usecase/config"
	"file_tool/domain/usecase/core_lib"
)

// Container manages and stores the lifetime of long-lived object instances
type Container struct {
	CoreLibRepo repository.CoreLibRepository
	ConfigRepo repository.ConfigRepository
	SetConfigUseCase config.SetConfigUseCase
	GetConfigUseCase config.GetConfigUseCase
	GetAllConfigUseCase config.GetAllConfigUseCase
	WriteCoreLibsUseCase core_lib.WriteCoreLibsUseCase
}

// NewContainer acts as the manual DI Initializer (Injector)
func NewContainer() *Container {

	configRepoInstance := repository_impl.NewConfigRepository()
	
	setConfigUseCaseInstance := config.NewSetConfigUseCase(configRepoInstance)
	getConfigUseCaseInstance := config.NewGetConfigUseCase(configRepoInstance)

	coreLibRepoInstance := repository_impl.NewCoreLibRepository(getConfigUseCaseInstance)

	return &Container{
		CoreLibRepo: coreLibRepoInstance,
		ConfigRepo: configRepoInstance,
		SetConfigUseCase: setConfigUseCaseInstance,
		GetConfigUseCase: getConfigUseCaseInstance,
		GetAllConfigUseCase: config.NewGetAllConfigUseCase(configRepoInstance),
		WriteCoreLibsUseCase: core_lib.NewWriteCoreLibsUseCase(coreLibRepoInstance),
	}
}