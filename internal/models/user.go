package models

import "gopkg.in/guregu/null.v3"

// represent user model
type User struct {
	ID           int64       `json:"id"; redis:"id"`
	Username     string      `json:"username"; redis:"username"`
	Password     string      `json:"password"; redis:"password"`
	Nickname     null.String `json:"nickname"; redis:"nickname"`
	ProfileImage null.String `json:"profile_image"; redis:"profile_image"`
}

type UserProfile struct {
	ID           int64       `json:"id"; redis:"id"`
	Username     string      `json:"username"; redis:"username"`
	Nickname     null.String `json:"nickname"; redis:"nickname"`
	ProfileImage null.String `json:"profile_image"; redis:"profile_image"`
}
