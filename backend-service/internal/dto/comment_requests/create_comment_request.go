package commentrequests

type CreateComment struct {
	Content string `json:"content" binding:"required"`
}
