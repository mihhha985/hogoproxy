package auth

type User struct {
	Username string `json:"username" example:"user"`
	Password string `json:"password" example:"password123"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type TokenClaims map[string]interface{}
