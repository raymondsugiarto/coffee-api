package entity

import "time"

type ReportParticipantSummary struct {
	ParticipantType     string `json:"participantType"`
	ParticipantCategory string `json:"participantCategory"`
	ParticipantCount    int64  `json:"participantCount"`
	PilarType           string `json:"pilarType"`
	DataSource          string `json:"dataSource"`
}

type ReportParticipantSummaryFilter struct {
	EndDate *time.Time
}

type SummaryType string

const (
	SummaryTypeAll               SummaryType = "ALL"
	SummaryTypePilar             SummaryType = "PILAR"
	SummaryTypeNonPilar          SummaryType = "NON PILAR"
	SummaryTypeManfaatLain       SummaryType = "MANFAAT LAIN"
	SummaryTypePesertaBadanUsaha SummaryType = "PESERTA BADAN USAHA"
)

type GroupedSummary struct {
	Type   SummaryType                `json:"pilarType"`
	Groups []ReportParticipantSummary `json:"groups"`
}
