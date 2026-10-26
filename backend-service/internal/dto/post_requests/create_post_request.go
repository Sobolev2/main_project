package postrequests

type CreatePostRequest struct {
	Content string `json:"content" binding:"required"`
}
