package entity

func NewUser(uid string, username string, email string, userType string) (
	*Users, error) {

	user := &Users{
		Uid:      uid,
		Username: username,
		Email:    email,
		Type:     userType,
	}

	err := user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Users) Validate() error {
	if u.Username == "" {
		return ErrInvalidUserNameBlank
	}

	if u.Email == "" {
		return ErrInvalidUserEmailBlank
	}

	if u.Uid == "" {
		return ErrInvalidUserUidBlank
	}

	if u.Type == "" {
		return ErrInvalidUserTypeBlank
	}

	return nil
}
