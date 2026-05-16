package provider

import (
	"file_tool/data/repository_impl"
	"file_tool/domain/repository"
	"file_tool/domain/usecase/config"
	"file_tool/domain/usecase/core_lib"
	"file_tool/domain/usecase/framework"
)

// Container manages and stores the lifetime of long-lived object instances
type Container struct {
	CoreLibRepo repository.CoreLibRepository
	ConfigRepo repository.ConfigRepository
	GenCleanFilesRepo repository.GenCleanFilesRepository
	SetConfigUseCase config.SetConfigUseCase
	GetConfigUseCase config.GetConfigUseCase
	GetAllConfigUseCase config.GetAllConfigUseCase
	WriteCoreLibsUseCase core_lib.WriteCoreLibsUseCase
	GenCleanDirsUseCase framework.GenCleanDirsUseCase
	GenCleanFilesUseCase framework.GenCleanFilesUseCase
}

// NewContainer acts as the manual DI Initializer (Injector)
func NewContainer() *Container {

	configRepoInstance := repository_impl.NewConfigRepository()

	
	setConfigUseCaseInstance := config.NewSetConfigUseCase(configRepoInstance)
	getConfigUseCaseInstance := config.NewGetConfigUseCase(configRepoInstance)

	coreLibRepoInstance := repository_impl.NewCoreLibRepository(getConfigUseCaseInstance)
	genCleanFilesRepoInstance := repository_impl.NewGenCleanFilesRepository(getConfigUseCaseInstance)
	return &Container{
		CoreLibRepo: coreLibRepoInstance,
		ConfigRepo: configRepoInstance,
		GenCleanFilesRepo: genCleanFilesRepoInstance,
		SetConfigUseCase: setConfigUseCaseInstance,
		GetConfigUseCase: getConfigUseCaseInstance,
		GetAllConfigUseCase: config.NewGetAllConfigUseCase(configRepoInstance),
		WriteCoreLibsUseCase: core_lib.NewWriteCoreLibsUseCase(coreLibRepoInstance),
		GenCleanDirsUseCase: framework.NewGenCleanDirsUseCase(genCleanFilesRepoInstance),
		GenCleanFilesUseCase: framework.NewGenCleanFilesUseCase(genCleanFilesRepoInstance),
	}
}