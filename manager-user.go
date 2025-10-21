package nuxeo

import (
	"context"
	"fmt"
	"log/slog"
)

type UserManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

func (um *UserManager) CurrentUser(ctx context.Context) (*User, error) {
	// first get the username via the login operation
	loginInfo := &struct {
		Username string `json:"username"`
	}{}
	if err := um.client.OperationManager().NewOperation(ctx, "login", &NuxeoRequestOptions{
		Enrichers: map[string][]string{
			"user": {"userprofile"},
		},
	}).ExecuteInto(loginInfo); err != nil {
		um.logger.Error("Failed to get current user", "error", err)
		return nil, err
	}

	// then fetch the user details
	return um.FetchUser(ctx, loginInfo.Username, nil)
}

func (um *UserManager) FetchUser(ctx context.Context, id string, options *NuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx).SetNuxeoOption(options).SetResult(&User{}).Get("/api/v1/user/" + id)

	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to fetch user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch user: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*User), nil
}

func (um *UserManager) CreateUser(ctx context.Context, user *User) (*User, error) {
	res, err := um.client.NewRequest(ctx).SetBody(User{
		EntityType: "user",
		Properties: user.Properties,
	}).SetResult(&User{}).Post("/api/v1/user/")

	if err != nil || res.StatusCode() != 201 {
		um.logger.Error("Failed to create user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to create user: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*User), nil
}

func (c *UserManager) UpdateUser(user *User) error {
	return nil
}
