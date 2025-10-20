package nuxeo

import (
	"context"
	"fmt"
)

// User represents a Nuxeo user.
type User struct {
	EntityType string         `json:"entity-type"`
	Id         string         `json:"id"`
	Properties map[string]any `json:"properties"`
}

func (c *NuxeoClient) FetchUser(ctx context.Context, id string, options *NuxeoRequestOption) (*User, error) {
	user := &User{}
	res, err := c.NewRequest(ctx).SetNuxeoOption(options).SetResult(user).Get("/api/v1/user/" + id)

	if err != nil || res.StatusCode() != 200 {
		c.logger.Error("Failed to fetch user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch user: %d %w", res.StatusCode(), err)
	}
	return user, nil
}

func (c *NuxeoClient) CreateUser(ctx context.Context, user *User) (*User, error) {
	createdUser := &User{}
	res, err := c.NewRequest(ctx).SetBody(User{
		EntityType: "user",
		Properties: user.Properties,
	}).SetResult(createdUser).Post("/api/v1/user/")

	if err != nil || res.StatusCode() != 201 {
		c.logger.Error("Failed to create user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to create user: %d %w", res.StatusCode(), err)
	}
	return createdUser, nil
}

func (c *NuxeoClient) UpdateUser(user *User) error {
	return nil
}
