package events

type PostLiked struct {
	Event  string `json:"event"`
	PostID int    `json:"post_id"`
	UserID int    `json:"user_id"`
}
