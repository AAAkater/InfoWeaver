package response

type TokenItem struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
