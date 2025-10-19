package nuxeo

// User represents a Nuxeo user.
type User struct {
	Username   string
	Email      string
	Properties map[string]any
}

// NewUser creates a new User instance.
func NewUser(username, email string, props map[string]any) *User {
	return &User{
		Username:   username,
		Email:      email,
		Properties: props,
	}
}
