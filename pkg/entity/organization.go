package entity

const (
	OriginKey       = "x-origin"
	OriginTypeKey   = "x-origin-type" // ADMIN, COMPANY, CUSTOMER
	OrganizationKey = "organization"

	UserContextKey        = "user"
	UserCredentialDataKey = "userCredentialData"

	CompanyKey = "company"
)

type OrganizationData struct {
	ID string `json:"id"`
}

type UserCredentialData struct {
	ID         string `json:"id"`   // user credential id
	UserID     string `json:"uid"`  // user id
	CustomerID string `json:"cid"`  // user id
	AccountID  string `json:"aid"`  // user id
	CompanyID  string `json:"coid"` // user id
}
