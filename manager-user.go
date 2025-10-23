package nuxeo

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
)

type UserManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

////////////////
//// GROUPS ////
////////////////

func (um *UserManager) FetchGroup(ctx context.Context, name string, options *nuxeoRequestOptions) (*Group, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).Get(apiV1 + "/group/" + name)

	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to fetch group", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch group: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) CreateGroup(ctx context.Context, group Group, options *nuxeoRequestOptions) (*Group, error) {
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).Post(apiV1 + "/group")

	if err != nil || res.StatusCode() != 201 {
		um.logger.Error("Failed to create group", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to create group: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) UpdateGroup(ctx context.Context, name string, group Group, options *nuxeoRequestOptions) (*Group, error) {
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).Put(apiV1 + "/group/" + name)

	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to update group", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to update group: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) DeleteGroup(ctx context.Context, name string, options *nuxeoRequestOptions) error {
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).Delete(apiV1 + "/group/" + name)

	if err != nil || res.StatusCode() != 204 {
		um.logger.Error("Failed to delete group", "error", err, "status", res.StatusCode())
		return fmt.Errorf("failed to delete group: %d %w", res.StatusCode(), err)
	}
	return nil
}

func (um *UserManager) SearchGroup(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
	reqQuery := "q=" + url.QueryEscape(query)
	if paginationQuery := paginationOptions.QueryParams(); paginationQuery != "" {
		reqQuery += "&" + paginationQuery
	}

	res, err := um.client.NewRequest(ctx, options).SetResult(&Groups{}).Get(apiV1 + "/group/search?" + reqQuery)

	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to search groups", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to search groups: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Groups), nil
}

func (um *UserManager) AttachGroupToUser(ctx context.Context, idOrGroupName string, idOrUsername string, options *nuxeoRequestOptions) (*Group, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).Post(apiV1 + "/group/" + idOrGroupName + "/user/" + idOrUsername)

	if err != nil || res.StatusCode() != 201 {
		um.logger.Error("Failed to attach group to user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to attach group to user: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) fetchGroupMemberUsers(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
	return nil, nil
}

func (um *UserManager) fetchGroupMemberGroups(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
	return nil, nil
}

///////////////
//// USERS ////
///////////////

func (um *UserManager) FetchUser(ctx context.Context, id string, options *nuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&User{}).Get(apiV1 + "/user/" + id)

	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to fetch user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch user: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*User), nil
}

func (um *UserManager) CreateUser(ctx context.Context, user *User, options *nuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx, options).SetBody(User{
		entity: entity{
			EntityType: "user",
		},
		Properties: user.Properties,
	}).SetResult(&User{}).Post(apiV1 + "/user")

	if err != nil || res.StatusCode() != 201 {
		um.logger.Error("Failed to create user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to create user: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*User), nil
}

func (c *UserManager) UpdateUser(ctx context.Context, idOrUsername string, user *User, options *nuxeoRequestOptions) (*User, error) {
	return nil, nil
}

func (c *UserManager) DeleteUser(ctx context.Context, idOrUsername string, options *nuxeoRequestOptions) error {
	return nil
}

func (c *UserManager) SearchUser(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
	return nil, nil
}

func (c *UserManager) AddUserToGroup(ctx context.Context, idOrUsername string, idOrGroupName string, options *nuxeoRequestOptions) (*User, error) {
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
	operationMgr := um.client.OperationManager()
	operationLogin := operationMgr.NewOperation("login")
	if err := operationMgr.ExecuteInto(ctx, *operationLogin, nil, loginInfo); err != nil {
		um.logger.Error("Failed to get current user", "error", err)
		return nil, err
	}

	// then fetch the user details
	return um.FetchUser(ctx, loginInfo.Username, nil)
}

func (um *UserManager) FetchWorkflowInstances(ctx context.Context, options *nuxeoRequestOptions) (*Workflows, error) {
	return nil, nil
}

func (um *UserManager) StartWorkflowInstance(ctx context.Context, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	return nil, nil
}
