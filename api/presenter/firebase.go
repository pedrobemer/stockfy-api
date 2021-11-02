package presenter

type SignUpBody struct {
	Password    string `json:"password,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type SignInBody struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type ForgotPasswordBody struct {
	Email string `json:"email,omitempty"`
}

type UserApiReturn struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type UserLoginApiReturn struct {
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	Expiration   string `json:"expiration"`
}

func ConvertUserToUserApiReturn(email string, displayName string) UserApiReturn {
	return UserApiReturn{
		Email:       email,
		DisplayName: displayName,
	}
}

func ConvertUserLoginToUserLoginApiReturn(email string, displayName string,
	idToken string, refreshToken string, expiration string) UserLoginApiReturn {
	return UserLoginApiReturn{
		Email:        email,
		DisplayName:  displayName,
		IdToken:      idToken,
		RefreshToken: refreshToken,
		Expiration:   expiration,
	}
}
