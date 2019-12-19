package models

import "gopkg.in/guregu/null.v3"

// represent user model
type User struct {
	ID           int64 `json:"id"`
	Username     string `json:"username"`
	Password	string `json:"password"`
	Nickname     null.String `json:"nickname"`
	ProfileImage null.String `json:"profile_image"`
}

type UserProfile struct {
	ID           int64      `json:id"`
	Username     string      `json:"username"`
	Nickname     null.String `json:"nickname"`
	ProfileImage null.String `json:"profile_image"`
}
