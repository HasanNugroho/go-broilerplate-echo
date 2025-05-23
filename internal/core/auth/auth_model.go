package auth

type AuthModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	Data         interface{} `json:"data"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"your-refresh-token"`
}
