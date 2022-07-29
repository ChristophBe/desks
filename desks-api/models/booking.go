package models

import (
	"reflect"
	"time"
)

type Booking struct {
	AbstractModel
	UserId int64     `json:"-"`
	User   User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RoomId int64     `json:"-"`
	Room   Room      `json:"room" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

func (b Booking) Patch(other Booking) Booking {
	if !reflect.ValueOf(other.Room).IsZero() {
		b.Room = other.Room
	}
	if !reflect.ValueOf(other.Start).IsZero() {
		b.Start = other.Start
	}
	if !reflect.ValueOf(other.End).IsZero() {
		b.End = other.End
	}
	return b

}
