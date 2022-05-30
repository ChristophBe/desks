package models

type User struct {
	AbstractModel
	ExternalUserId string `json:"-" gorm:"unique;notNull"`
	FamilyName     string `json:"familyName" gorm:"notNull"`
	GivenName      string `json:"givenName" gorm:"notNull"`
}
