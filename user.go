package contentful

// UsersService service
type UsersService service

// User model
type User struct {
	Sys                            *Sys   `json:"sys,omitempty"`
	FirstName                      string `json:"firstName"`
	LastName                       string `json:"lastName"`
	AvatarURL                      string `json:"avatarUrl"`
	Email                          string `json:"email"`
	Activated                      bool   `json:"activated"`
	SignInCount                    int    `json:"signInCount"`
	Confirmed                      bool   `json:"confirmed"`
	TwoFactorAuthenticationEnabled bool   `json:"2faEnabled"`
}

// Me returns current authenticated user
func (service *UsersService) Me() (*User, error) {
	path := "/users/me"
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := service.c.do(req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
