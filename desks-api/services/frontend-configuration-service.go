package services

import (
	"github.com/ChristophBe/desks/desks-api/dtos"
	"github.com/ChristophBe/desks/desks-api/util"
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
	maxDataAge := util.GetIntEnvironmentVariable("MAX_USERDATA_AGE_DAYS", 90) * 24
	return dtos.FrontendConfigurationDTO{
		MaximalBookingDate: time.Now().Round(24 * time.Hour).Add(time.Duration(maxDataAge) * time.Hour),
	}
}
