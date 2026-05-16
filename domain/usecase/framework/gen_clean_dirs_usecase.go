package framework

import (
	"file_tool/domain/repository"
	"file_tool/domain/entity"
)

type GenCleanDirsUseCase interface {
	GenCleanDirs(frame_work entity.Framework, path string) (string, error)
}

type genCleanDirsUseCase struct {
	repo repository.GenCleanFilesRepository
}

func NewGenCleanDirsUseCase(genCleanFilesRepo repository.GenCleanFilesRepository) GenCleanDirsUseCase {
	return &genCleanDirsUseCase{
		repo: genCleanFilesRepo,
	}
}

func (uc *genCleanDirsUseCase) GenCleanDirs(frame_work entity.Framework, path string) (string, error) {
	return uc.repo.GenCleanDirs(frame_work, path)
}