package repository

import (
	"file_tool/domain/entity"
)

type GenCleanFilesRepository interface {
	GenCleanDirs(frame_work entity.Framework, path string) (string, error)
	GenCleanFiles(frame_work entity.Framework, path string, name string) (error)
}
