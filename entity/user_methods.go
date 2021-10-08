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
		return ErrInvalidUserName
	}

	if u.Email == "" {
		return ErrInvalidUserEmail
	}

	if u.Uid == "" {
		return ErrInvalidUserUid
	}

	if u.Type == "" {
		return ErrInvalidUserType
	}

	return nil
}
