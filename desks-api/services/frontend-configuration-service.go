package services

import (
	"github.com/ChristophBe/desks/desks-api/dtos"
	"github.com/ChristophBe/desks/desks-api/util/configuration"
	"time"
)

type FrontendConfigurationService interface {
	GetFrontendConfiguration() dtos.FrontendConfigurationDTO
}

func NewFrontendConfigurationService() FrontendConfigurationService {
	return new(frontendConfigurationServiceImpl)
}

type frontendConfigurationServiceImpl struct{}

func (f frontendConfigurationServiceImpl) GetFrontendConfiguration() dtos.FrontendConfigurationDTO {
	maxDataAge := configuration.MaxUserdataAgeDays.GetValue() * 24
	return dtos.FrontendConfigurationDTO{
		MaximalBookingDate: time.Now().Round(24 * time.Hour).Add(time.Duration(maxDataAge) * time.Hour),
	}
}
