package framework

import (
	"file_tool/domain/entity"
	"file_tool/domain/repository"
)

type GenCleanFilesUseCase interface {
	GenCleanFiles(frame_work entity.Framework, path string, name string) (error)
}

type genCleanFilesUseCase struct {
	repo repository.GenCleanFilesRepository
}

func NewGenCleanFilesUseCase(genCleanFilesRepo repository.GenCleanFilesRepository) GenCleanFilesUseCase {
	return &genCleanFilesUseCase{
		repo: genCleanFilesRepo,
	}
}

func (uc *genCleanFilesUseCase) GenCleanFiles(frame_work entity.Framework, path string, name string) (error) {
	return uc.repo.GenCleanFiles(frame_work, path, name)
}