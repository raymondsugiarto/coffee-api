package entity

import "github.com/raymondsugiarto/coffee-api/pkg/model"

type UserCredentialDto struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organizationId"`
	Username       string   `json:"username"`
	Password       string   `json:"password"`
	Email          string   `json:"email" `
	CustomerID     string   `json:"customerId"`
	User           *UserDto `json:"user"`
}

func (u *UserCredentialDto) ToModel() *model.UserCredential {
	m := &model.UserCredential{
		OrganizationID: u.OrganizationID,
		Username:       u.Username,
		Password:       u.Password,
	}
	if u.ID != "" {
		m.ID = u.ID
	}
	return m
}

func (u *UserCredentialDto) FromModel(m *model.UserCredential) *UserCredentialDto {
	return &UserCredentialDto{
		ID:       m.ID,
		Username: m.Username,
		Password: m.Password,
		User:     new(UserDto).FromModel(&m.User),
	}
}

type ChangePasswordDto struct {
	UserCredentialID string `json:"userCredentialId"`
	Password         string `json:"password"`
}

type PasswordChangeInputDto struct {
	UserID          string
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
