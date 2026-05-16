package core_lib

import (
	"file_tool/domain/repository"
)

type WriteCoreLibsUseCase interface {
	WriteCoreLibs() bool
}

type writeCoreLibsUseCase struct {
	coreLibRepo repository.CoreLibRepository
}

func NewWriteCoreLibsUseCase(coreLibRepo repository.CoreLibRepository) WriteCoreLibsUseCase {
	return &writeCoreLibsUseCase{
		coreLibRepo: coreLibRepo,
	}
}

func (uc *writeCoreLibsUseCase) WriteCoreLibs() bool {
	return uc.coreLibRepo.WriteCoreLibs()
}