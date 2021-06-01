package aggregate

type User struct {
	ID          string
	Username    string
	DisplayName string
}

func NewUser(id, username, displayName string) *User {
	return &User{
		ID:          id,
		Username:    username,
		DisplayName: displayName,
	}
}
