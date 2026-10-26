package commentrequests

type UpdateComment struct {
	Content string `json:"content" binding:"required"`
}
