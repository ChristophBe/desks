package dtos

import (
	"github.com/ChristophBe/desks/desks-api/models"
	"time"
)

type BookingDefaultsDto struct {
	Room  models.Room `json:"room"`
	Start time.Time   `json:"start"`
	End   time.Time   `json:"end"`
}
