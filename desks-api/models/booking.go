package models

import "time"

type Booking struct {
	AbstractModel
	UserId int64     `json:"-"`
	User   User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RoomId int64     `json:"-"`
	Room   Room      `json:"room" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}
