package presenter

type SignUpBody struct {
	Password    string `json:"password,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type ForgotPasswordBody struct {
	Email string `json:"email,omitempty"`
}

type UserApiReturn struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

func ConvertUserToUserApiReturn(email string, displayName string) UserApiReturn {
	return UserApiReturn{
		Email:       email,
		DisplayName: displayName,
	}
}
