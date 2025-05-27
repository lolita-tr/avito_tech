package auth

type AuthRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Token string `json:"token"`
}
