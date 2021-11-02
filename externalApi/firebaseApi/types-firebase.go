package firebaseApi

type EmailVerificationParams struct {
	RequestType string `json:"requestType"`
	IdToken     string `json:"idToken"`
}

type PasswordReset struct {
	RequestType string `json:"requestType,omitempty"`
	Email       string `json:"email,omitempty"`
}

type UserLogin struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}
