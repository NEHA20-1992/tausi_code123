package model

type AuthenticationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      string `json:"otp,omitempty"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email"`
	ResetCode       string `json:"resetCode"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type RefreshAuthenticationRequest struct {
	RefreshToken string `json:"refreshToken"`
}
