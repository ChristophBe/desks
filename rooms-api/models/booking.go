package models

import "time"

type Booking struct {
	AbstractModel
	UserId int64     `json:"-"`
	User   User      `json:"user"`
	RoomId int64     `json:"-"`
	Room   Room      `json:"room"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}
