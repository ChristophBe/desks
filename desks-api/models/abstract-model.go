package models

import "time"

type AbstractModel struct {
	Id      int64     `gorm:"primaryKey" json:"id"`
	Created time.Time `gorm:"autoCreateTime" json:"created"`
	Updated time.Time `gorm:"autoUpdateTime" json:"updated"`
}
