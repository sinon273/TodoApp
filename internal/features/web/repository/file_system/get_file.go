package web_fs_repository

import (
	core_error "TodoApp/internal/core/errors"
	web_domain "TodoApp/internal/features/web/domain"
	"fmt"
	"os"
)

func (r *WebRepository) GetFile(filePath string) (web_domain.File, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return web_domain.File{}, fmt.Errorf(
				"file: %s: %w",
				filePath,
				core_error.ErrNotFound,
			)
		}
		return web_domain.File{}, fmt.Errorf(
			"get file: %s: %w",
			filePath,
			err,
		)
	}

	return web_domain.NewFile(file, filePath), nil
}
