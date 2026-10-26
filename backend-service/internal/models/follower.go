package models

type Follower struct {
	ID         int `json:"id"`
	FollowerID int `json:"follower_id"`
	UserID     int `json:"user_id"`
}
