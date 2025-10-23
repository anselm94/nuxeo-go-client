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

func (um *UserManager) FetchGroupMemberUsers(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
	urlPath := apiV1 + "/group/" + idOrGroupName + "/@users"
	if paginationOptions != nil {
		if paginationQuery := paginationOptions.QueryParams(); paginationQuery != "" {
			urlPath += "?" + paginationQuery
		}
	}
	res, err := um.client.NewRequest(ctx, options).SetResult(&Users{}).Get(urlPath)
	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to fetch group member users", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch group member users: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Users), nil
}

func (um *UserManager) FetchGroupMemberGroups(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
	urlPath := apiV1 + "/group/" + idOrGroupName + "/@groups"
	if paginationOptions != nil {
		if paginationQuery := paginationOptions.QueryParams(); paginationQuery != "" {
			urlPath += "?" + paginationQuery
		}
	}
	res, err := um.client.NewRequest(ctx, options).SetResult(&Groups{}).Get(urlPath)
	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to fetch group member groups", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch group member groups: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Groups), nil
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

func (um *UserManager) CreateUser(ctx context.Context, user User, options *nuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx, options).SetBody(User{
		entity: entity{
			EntityType: EntityTypeUser,
		},
		Properties: user.Properties,
	}).SetResult(&User{}).Post(apiV1 + "/user")

	if err != nil || res.StatusCode() != 201 {
		um.logger.Error("Failed to create user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to create user: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*User), nil
}

func (um *UserManager) UpdateUser(ctx context.Context, idOrUsername string, user User, options *nuxeoRequestOptions) (*User, error) {
	// Ensure entity-type and id are set
	updatePayload := &User{
		entity: entity{
			EntityType: EntityTypeUser,
		},
		Id:         idOrUsername,
		Properties: user.Properties,
	}
	res, err := um.client.NewRequest(ctx, options).SetBody(updatePayload).SetResult(&User{}).Put(apiV1 + "/user/" + idOrUsername)
	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to update user", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to update user: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*User), nil
}

func (um *UserManager) DeleteUser(ctx context.Context, idOrUsername string, options *nuxeoRequestOptions) error {
	res, err := um.client.NewRequest(ctx, options).Delete(apiV1 + "/user/" + idOrUsername)
	if err != nil || res.StatusCode() != 204 {
		um.logger.Error("Failed to delete user", "error", err, "status", res.StatusCode())
		return fmt.Errorf("failed to delete user: %d %w", res.StatusCode(), err)
	}
	return nil
}

func (um *UserManager) SearchUsers(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
	reqQuery := "q=" + url.QueryEscape(query)
	if paginationOptions != nil {
		if paginationQuery := paginationOptions.QueryParams(); paginationQuery != "" {
			reqQuery += "&" + paginationQuery
		}
	}
	res, err := um.client.NewRequest(ctx, options).SetResult(&Users{}).Get(apiV1 + "/user/search?" + reqQuery)
	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to search users", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to search users: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Users), nil
}

func (um *UserManager) AddUserToGroup(ctx context.Context, idOrUsername string, idOrGroupName string, options *nuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&User{}).Post(apiV1 + "/user/" + idOrUsername + "/group/" + idOrGroupName)
	if err != nil || res.StatusCode() != 201 {
		um.logger.Error("Failed to add user to group", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to add user to group: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*User), nil
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
	res, err := um.client.NewRequest(ctx, options).SetResult(&Workflows{}).Get(apiV1 + "/workflow")
	if err != nil || res.StatusCode() != 200 {
		um.logger.Error("Failed to fetch workflow instances", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch workflow instances: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Workflows), nil
}

func (um *UserManager) StartWorkflowInstance(ctx context.Context, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	// Ensure entity-type is set
	workflow.entity.EntityType = "workflow"
	res, err := um.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&Workflow{}).Post(apiV1 + "/workflow")
	if err != nil || res.StatusCode() != 201 {
		um.logger.Error("Failed to start workflow instance", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to start workflow instance: %d %w", res.StatusCode(), err)
	}
	return res.Result().(*Workflow), nil
}
