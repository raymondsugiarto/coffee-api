package entity

type ReportAUMData struct {
	CompanyName string  `json:"companyName"`
	PilarType   string  `json:"pilarType"`
	TotalAUM    float64 `json:"totalAum"`
}

type ReportAUMCompanyType string

const (
	ReportAUMCompanyTypeDKP  ReportAUMCompanyType = "DKP"
	ReportAUMCompanyTypePPIP ReportAUMCompanyType = "PPIP"
)

type ReportAUMFilter struct {
	Month       int                  `json:"month" validate:"required,min=1,max=12"`
	Year        int                  `json:"year" validate:"required,min=1900"`
	CompanyType ReportAUMCompanyType `json:"companyType" validate:"required,oneof=DKP PPIP"`
}

type GroupedAum struct {
	Type   SummaryType     `json:"type"`
	Groups []ReportAUMData `json:"groups"`
}
