package dto

type StartAuthRequest struct {
	Email string `json:"email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LogoutRequest struct {
	Email string `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
}
