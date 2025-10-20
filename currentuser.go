package nuxeo

import (
	"context"
)

type resLoginInfo struct {
	Username string `json:"username"`
}

func (c *NuxeoClient) CurrentUser(ctx context.Context) (*User, error) {
	// first get the username via the login operation
	loginInfo := &resLoginInfo{}
	if err := c.NewOperation(ctx, "login", &NuxeoRequestOption{
		Enrichers: map[string][]string{
			"user": {"userprofile"},
		},
	}).ExecuteInto(loginInfo); err != nil {
		c.logger.Error("Failed to get current user", "error", err)
		return nil, err
	}

	// then fetch the user details
	return c.FetchUser(ctx, loginInfo.Username, nil)
}
