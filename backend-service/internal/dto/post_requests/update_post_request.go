package postrequests

type UpdatePostRequest struct {
	Content string `json:"content" binding:"required"`
}
