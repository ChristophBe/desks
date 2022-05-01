package models

type User struct {
	AbstractModel
	Username  string `json:"username" gorm:"unique;notNull"`
	Firstname string `json:"firstname" gorm:"notNull"`
	Lastname  string `json:"lastname" gorm:"notNull"`
}
