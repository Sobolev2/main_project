package dto

type UpdateUserRequest struct {
	UserName  string `json:"user_name" `
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
