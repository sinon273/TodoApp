package web_service

import web_domain "TodoApp/internal/features/web/domain"

type WebService struct {
	webRepository WebRepository
}

func NewWebService(webRepository WebRepository) *WebService {
	return &WebService{webRepository: webRepository}
}

type WebRepository interface {
	GetFile(filePath string) (web_domain.File, error)
}
