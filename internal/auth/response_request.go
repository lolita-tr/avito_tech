package auth

type AuthRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
