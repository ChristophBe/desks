package dtos

import "time"

type FrontendConfigurationDTO struct {
	MaximalBookingDate time.Time `json:"maximalBookingDate"`
}
