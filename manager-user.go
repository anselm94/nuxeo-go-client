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

////////////////
//// GROUPS ////
////////////////

func (um *UserManager) FetchGroup(ctx context.Context, name string) (*Group, error) {
	return nil, nil
}

func (um *UserManager) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	return nil, nil
}

func (um *UserManager) UpdateGroup(ctx context.Context, name string, group *Group) (*Group, error) {
	return nil, nil
}

func (um *UserManager) DeleteGroup(ctx context.Context, name string) error {
	return nil
}

func (um *UserManager) SearchGroup(ctx context.Context, query string, paginationOptions *PaginationOptions) (*Groups, error) {
	return nil, nil
}

func (um *UserManager) AttachGroupToUser(ctx context.Context, idOrGroupName string, idOrUsername string) (*Group, error) {
	return nil, nil
}

func (um *UserManager) fetchGroupMemberUsers(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions) (*Users, error) {
	return nil, nil
}

func (um *UserManager) fetchGroupMemberGroups(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions) (*Groups, error) {
	return nil, nil
}

///////////////
//// USERS ////
///////////////

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

func (c *UserManager) UpdateUser(ctx context.Context, idOrUsername string, user *User) (*User, error) {
	return nil, nil
}

func (c *UserManager) DeleteUser(ctx context.Context, idOrUsername string) error {
	return nil
}

func (c *UserManager) SearchUser(ctx context.Context, query string, paginationOptions *PaginationOptions) (*Users, error) {
	return nil, nil
}

func (c *UserManager) AddUserToGroup(ctx context.Context, idOrUsername string, idOrGroupName string) (*User, error) {
	return nil, nil
}

//////////////////////
//// CURRENT USER ////
//////////////////////

func (um *UserManager) FetchCurrentUser(ctx context.Context) (*User, error) {
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

func (um *UserManager) FetchWorkflowInstances(ctx context.Context) (*Workflows, error) {
	return nil, nil
}

func (um *UserManager) StartWorkflowInstance(ctx context.Context, workflow Workflow) (*Workflow, error) {
	return nil, nil
}
