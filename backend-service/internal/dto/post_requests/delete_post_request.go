package postrequests

type DeletePostRequest struct {
	PostID int `json:"post_id" binding:"required"`
}
