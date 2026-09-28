package users_transport_http

import (
	"TodoApp/internal/core/domain"

	"github.com/google/uuid"
)

type UserDTOResponse struct {
	ID          uuid.UUID `json:"id" example:"f25b2127-1f7a-40c8-bbca-d3c766e99565"`
	Version     int       `json:"version" example:"3"`
	FullName    string    `json:"full_name" example:"Ivan Ivanov"`
	PhoneNumber *string   `json:"phone_number" example:"78005353535"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	userDTO := make([]UserDTOResponse, len(users))
	for i, user := range users {
		userDTO[i] = userDTOFromDomain(user)
	}

	return userDTO
}
