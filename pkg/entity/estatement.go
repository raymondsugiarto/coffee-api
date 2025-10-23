package entity

import "time"

type EstatementRequestDto struct {
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	CustomerID string    `json:"customerId"`
}

type EstatementEmailRequestDto struct {
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	CustomerID string    `json:"customerId"`
}

type EstatementDto struct {
	FilePath string `json:"filePath"`
}
