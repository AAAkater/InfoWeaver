package response

type TokenResult struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type UserInfoResult struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
