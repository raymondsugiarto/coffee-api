package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Reward struct {
	concern.CommonWithIDs
	Name        string
	Points      int
	Stock       int
	ImageUrl    string
	ExpiredAt   time.Time
	Status      RewardStatus
	Description string
	Code        string
}

type RewardStatus string

const (
	REWARD_ACTIVE   RewardStatus = "ACTIVE"
	REWARD_INACTIVE RewardStatus = "INACTIVE"
)
