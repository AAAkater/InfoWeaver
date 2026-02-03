package models

type RegisterReq struct {
	Username string `json:"username" validate:"required,min=4,max=16"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=18"`
}

type UserLoginReq struct {
	Username string `json:"username" validate:"required,min=4,max=16"`
	Password string `json:"password" validate:"required,min=6,max=18"`
}
type UserLoginResp struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type ResetPasswordReq struct {
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateUserInfoReq struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type UserInfoResp struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
