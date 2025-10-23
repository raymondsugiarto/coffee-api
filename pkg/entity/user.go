package entity

import (
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type CustomerAccountListItem struct {
	ID                 string `json:"id" bson:"id"`
	AccountCode        string `json:"accountCode" bson:"accountCode"`
	UserID             string `json:"userId" bson:"userId"`
	CreatedAt          string `json:"createdAt" bson:"createdAt"`
	CustomerFollowerID string `json:"customerFollowerId" bson:"customerFollowerId"`
}

type MyAccountProfile struct {
	*UserCredentialData
	ID             string `json:"id" bson:"id"`
	AccountCode    string `json:"accountCode" bson:"accountCode"`
	Email          string `json:"email" bson:"email"`
	PhoneNumber    string `json:"phoneNumber" bson:"phoneNumber"`
	FirstName      string `json:"firstName" bson:"firstName"`
	LastName       string `json:"lastName" bson:"lastName"`
	UserID         string `json:"userId" bson:"userId"`
	CreatedAt      string `json:"createdAt" bson:"createdAt"`
	FollowerCount  int    `json:"followerCount" bson:"followerCount"`
	FollowingCount int    `json:"followingCount" bson:"followingCount"`
}

type CustomerAccount struct {
	*OrganizationData
	*UserCredentialData
	ID          string `json:"id" bson:"id"`
	AccountCode string `json:"accountCode" bson:"accountCode"`
	UserID      string `json:"userId" bson:"userId"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
}

type UserCredential struct {
	*OrganizationData
	*UserCredentialData
	ID         string   `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email" `
	CustomerID string   `json:"customerId"`
	User       *UserDto `json:"user"`
}

type CreateUser struct {
	*OrganizationData
	UserType    UserType `json:"userType"`
	UserID      string   `json:"userId"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phoneNumber"`
	Name        string   `json:"name"`
	Password    string   `json:"password"`
	Username    string   `json:"username"`
}

type UserDto struct {
	ID                      string               `json:"id"`
	OrganizationID          string               `json:"organizationId"`
	UserType                UserType             `json:"userType"`
	PhoneVerificationStatus model.IdentityStatus `json:"phoneVerificationStatus"`
	EmailVerificationStatus model.IdentityStatus `json:"emailVerificationStatus"`

	UserCredential []UserCredentialDto `json:"userCredential"`
	UserHasRoleDto []UserHasRoleDto    `json:"userHasRole"`
}

func (u *UserDto) ToModel() *model.User {
	m := &model.User{
		OrganizationID:          u.OrganizationID,
		UserType:                model.UserType(string(u.UserType)),
		PhoneVerificationStatus: u.PhoneVerificationStatus,
		EmailVerificationStatus: u.EmailVerificationStatus,
	}
	if u.ID != "" {
		m.ID = u.ID
	}

	if len(u.UserCredential) > 0 {
		m.UserCredential = make([]*model.UserCredential, 0)
		for _, uc := range u.UserCredential {
			m.UserCredential = append(m.UserCredential, uc.ToModel())
		}
	}
	if len(u.UserHasRoleDto) > 0 {
		m.UserHasRole = make([]*model.UserHasRole, 0)
		for _, hr := range u.UserHasRoleDto {
			m.UserHasRole = append(m.UserHasRole, hr.ToModel())
		}
	}
	return m
}

func (u *UserDto) FromModel(m *model.User) *UserDto {
	return &UserDto{
		OrganizationID:          u.OrganizationID,
		ID:                      m.ID,
		UserType:                UserType(m.UserType),
		PhoneVerificationStatus: m.PhoneVerificationStatus,
		EmailVerificationStatus: m.EmailVerificationStatus,
	}
}

const (
	CUSTOMER UserType = "CUSTOMER"
	ADMIN    UserType = "ADMIN"
	COMPANY  UserType = "COMPANY"
)

type UserType string

// Function to convert string to UserType
func StringToUserType(s string) (UserType, error) {
	switch s {
	case "ADMIN":
		return ADMIN, nil
	case "COMPANY":
		return COMPANY, nil
	case "CUSTOMER":
		return CUSTOMER, nil
	default:
		return "", errors.New("enumParseError")
	}
}
