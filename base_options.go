package nuxeo

// BaseOptions holds configuration for the NuxeoClient.
type BaseOptions struct {
	BaseURL  string
	User     string
	Password string
	Token    string
	Timeout  int // seconds
}

// WithBaseURL sets the base URL for the client.
func WithBaseURL(url string) Option {
	return func(c *NuxeoClient) {
		c.options.BaseURL = url
	}
}

// WithUser sets the username for authentication.
func WithUser(user string) Option {
	return func(c *NuxeoClient) {
		c.options.User = user
	}
}

// WithPassword sets the password for authentication.
func WithPassword(password string) Option {
	return func(c *NuxeoClient) {
		c.options.Password = password
	}
}

// WithToken sets the token for authentication.
func WithToken(token string) Option {
	return func(c *NuxeoClient) {
		c.options.Token = token
	}
}

// WithTimeout sets the timeout for requests.
func WithTimeout(timeout int) Option {
	return func(c *NuxeoClient) {
		c.options.Timeout = timeout
	}
}
