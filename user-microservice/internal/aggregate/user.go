package aggregate

type User struct {
	ID          string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

func NewUser(id, username, displayName string) *User {
	return &User{
		ID:          id,
		Username:    username,
		DisplayName: displayName,
	}
}
