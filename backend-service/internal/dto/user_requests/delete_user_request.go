package dto

type DeleteUserRequest struct {
	Password string `json:"password" binding:"required"`
}
