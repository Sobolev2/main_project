package dto

type LogoutAllDevicesExceptThisRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
