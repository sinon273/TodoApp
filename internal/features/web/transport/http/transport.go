package web_transport_http

import (
	core_http_server "TodoApp/internal/core/transport/http/server"
	web_domain "TodoApp/internal/features/web/domain"
)

type WebHTTPHandler struct {
	webService WebService
}

func NewWebHTTPHandler(webService WebService) *WebHTTPHandler {
	return &WebHTTPHandler{webService: webService}
}

type WebService interface {
	GetMainPage() (web_domain.File, error)
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
