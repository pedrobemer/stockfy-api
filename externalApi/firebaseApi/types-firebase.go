package firebaseApi

type EmailVerificationParams struct {
	RequestType string `json:"requestType"`
	IdToken     string `json:"idToken"`
}
