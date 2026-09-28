package web_service

import (
	web_domain "TodoApp/internal/features/web/domain"
	"fmt"
	"os"
	"path"
)

func (s *WebService) GetMainPage() (web_domain.File, error) {
	htmlFilePath := path.Join(
		os.Getenv("PROJECT_ROOT"),
		"public/index.html",
	)

	html, err := s.webRepository.GetFile(htmlFilePath)
	if err != nil {
		return web_domain.File{}, fmt.Errorf("get file from repository: %w", err)
	}

	return html, nil
}
