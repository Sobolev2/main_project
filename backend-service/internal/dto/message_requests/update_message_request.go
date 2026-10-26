package messagerequests

type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
