package dto

type UpdatePasswordRequest struct {
	LastPassword string `json:"last_password" binding:"required"`
	NewPassword  string `json:"new_password" binding:"required"`
}
