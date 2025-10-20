package nuxeo

import (
	"context"
	"encoding/json"
)

type resLoginInfo struct {
	Username string `json:"username"`
}

func (c *NuxeoClient) CurrentUser(ctx context.Context) (*User, error) {
	resOperation, err := c.NewOperation(ctx, "login", nil).Execute()
	if err != nil {
		c.logger.Error("Failed to get current user", "error", err)
		return nil, err
	}
	defer resOperation.Close()

	loginInfo := &resLoginInfo{}
	err = json.NewDecoder(resOperation).Decode(loginInfo)
	if err != nil {
		c.logger.Error("Failed to get current user", "error", err)
		return nil, err
	}

	return c.FetchUser(ctx, loginInfo.Username, nil)
}
