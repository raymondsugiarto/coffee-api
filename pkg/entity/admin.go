package entity

import "github.com/raymondsugiarto/coffee-api/pkg/model"

type AdminDto struct {
	ID              string   `json:"id"`
	AdminType       string   `json:"adminType"`
	UserID          string   `json:"userId"`
	User            *UserDto `json:"user"`
	PhoneNumber     string   `json:"phoneNumber"`
	Email           string   `json:"email"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	ProfileImageUrl string   `json:"profileImageUrl"`
	OrganizationID  string   `json:"organizationID,omitempty"`
	CompanyID       string   `json:"companyId,omitempty"`
}

type CreateAdminCompany struct {
	ID              string   `json:"id"`
	AdminType       string   `json:"adminType"`
	UserID          string   `json:"userId"`
	User            *UserDto `json:"user"`
	PhoneNumber     string   `json:"phoneNumber"`
	Email           string   `json:"email"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	ProfileImageUrl string   `json:"profileImageUrl"`
	OrganizationID  string   `json:"organizationID,omitempty"`
	CompanyID       string   `json:"companyId"`
	Password        string   `json:"password"`
}

func (dto *AdminDto) FromModel(m *model.Admin) *AdminDto {
	dto.ID = m.ID
	dto.AdminType = m.AdminType
	dto.UserID = m.UserID
	dto.PhoneNumber = m.PhoneNumber
	dto.Email = m.Email
	dto.FirstName = m.FirstName
	dto.LastName = m.LastName
	dto.ProfileImageUrl = m.ProfileImageUrl
	dto.OrganizationID = m.OrganizationID
	dto.CompanyID = m.CompanyID
	return dto
}

func (dto *AdminDto) ToModel() *model.Admin {
	m := &model.Admin{
		AdminType:       dto.AdminType,
		UserID:          dto.UserID,
		PhoneNumber:     dto.PhoneNumber,
		Email:           dto.Email,
		FirstName:       dto.FirstName,
		LastName:        dto.LastName,
		ProfileImageUrl: dto.ProfileImageUrl,
		OrganizationID:  dto.OrganizationID,
		CompanyID:       dto.CompanyID,
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	return m
}

func (dto *CreateAdminCompany) ToModel() *model.Admin {
	m := &model.Admin{
		AdminType:       dto.AdminType,
		UserID:          dto.UserID,
		PhoneNumber:     dto.PhoneNumber,
		Email:           dto.Email,
		FirstName:       dto.FirstName,
		LastName:        dto.LastName,
		ProfileImageUrl: dto.ProfileImageUrl,
		OrganizationID:  dto.OrganizationID,
		CompanyID:       dto.CompanyID,
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	if dto.User != nil {
		m.User = dto.User.ToModel()
	}
	return m
}

func (dto *CreateAdminCompany) FromModel(m *model.Admin) *CreateAdminCompany {
	dto.ID = m.ID
	dto.AdminType = m.AdminType
	dto.UserID = m.UserID
	dto.PhoneNumber = m.PhoneNumber
	dto.Email = m.Email
	dto.FirstName = m.FirstName
	dto.LastName = m.LastName
	dto.ProfileImageUrl = m.ProfileImageUrl
	dto.OrganizationID = m.User.OrganizationID
	return dto
}
