package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

// AdminCompany is the join table between admin and company — one row
// per (admin_id, company_id). Most callers issue raw JOIN queries
// against this table; this struct exists so we can read company_id
// for a given admin_id without raw SQL.
type AdminCompany struct {
	concern.CommonWithIDs
	OrganizationID string
	CompanyID      string
	AdminID        string
}
