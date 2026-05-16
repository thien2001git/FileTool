package repository

import (
	"file_tool/domain/entity"
)

type CoreLibRepository interface {
	// GetCurrentCoreLibs returns a list of currently installed core libraries.
	GetCurrentCoreLibs() ([]entity.CoreLib, error)
	// WriteCoreLibs writes the list of core libraries to storage.
	WriteCoreLibs() bool
}
