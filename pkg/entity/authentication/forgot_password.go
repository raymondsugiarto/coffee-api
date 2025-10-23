package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type ForgotPasswordInputDto struct {
	Email string `json:"email"`
}

func (f *ForgotPasswordInputDto) ToDto() *ForgotPasswordDto {
	return &ForgotPasswordDto{
		Email: f.Email,
	}
}

type ForgotPasswordDto struct {
	UserIdentityVerificationID string `json:"userIdentityVerificationId"`
	Email                      string `json:"email"`
}

func (f *ForgotPasswordDto) ToUserIdentityVerificationDto(identityFor model.IdentityFor) *UserIdentityVerificationDto {
	return &UserIdentityVerificationDto{
		UserIdentity: f.Email,
		IdentityFor:  string(identityFor),
		IdentityType: "EMAIL",
	}
}

type UserIdentityVerificationDto struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	IdentityFor    string    `json:"identityFor"`
	IdentityType   string    `json:"identityType"`
	UserIdentity   string    `json:"userIdentity"`
	UniqueCode     string    `json:"uniqueCode"`
	TryCount       int       `json:"tryCount"`
	ExpiredAt      time.Time `json:"expiredAt"`
	Status         string    `json:"status"`

	Data interface{} `json:"data"`
}

func (u *UserIdentityVerificationDto) ToModel() *model.UserIdentityVerification {
	m := &model.UserIdentityVerification{
		UserID:         u.UserID,
		OrganizationID: u.OrganizationID,
		IdentityFor:    model.IdentityFor(u.IdentityFor),
		IdentityType:   u.IdentityType,
		UserIdentity:   u.UserIdentity,
		UniqueCode:     u.UniqueCode,
		ExpiredAt:      u.ExpiredAt,
		TryCount:       u.TryCount,
		Status:         model.UserIdentityVerificationStatus(u.Status),
	}
	if u.ID != "" {
		m.ID = u.ID
	}
	return m
}

func (u *UserIdentityVerificationDto) FromModel(m *model.UserIdentityVerification) *UserIdentityVerificationDto {
	u.ID = m.ID
	u.UserID = m.UserID
	u.OrganizationID = m.OrganizationID
	u.IdentityFor = string(m.IdentityFor)
	u.IdentityType = m.IdentityType
	u.UserIdentity = m.UserIdentity
	u.UniqueCode = m.UniqueCode
	u.ExpiredAt = m.ExpiredAt.Local()
	u.TryCount = m.TryCount
	u.Status = string(m.Status)
	return u
}

type UserIdentityVerificationInputPasswordDto struct {
	ID         string `json:"id"`
	UniqueCode string `json:"uniqueCode"`
	Password   string `json:"password"`
}

func (u *UserIdentityVerificationInputPasswordDto) ToDto() *UserIdentityVerificationDto {
	return &UserIdentityVerificationDto{
		ID:         u.ID,
		UniqueCode: u.UniqueCode,
		Data:       u,
	}
}
