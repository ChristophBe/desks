package models

type Room struct {
	AbstractModel
	Name     string `json:"name" gorm:"unique;notNull"`
	Capacity int    `json:"capacity" gorm:"notNull"`
}
