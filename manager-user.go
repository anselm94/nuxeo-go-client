package nuxeo

import (
	"context"
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
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Get(apiV1 + "/group/" + name)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) CreateGroup(ctx context.Context, group Group, options *nuxeoRequestOptions) (*Group, error) {
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).SetError(&NuxeoError{}).Post(apiV1 + "/group")

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) UpdateGroup(ctx context.Context, name string, group Group, options *nuxeoRequestOptions) (*Group, error) {
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).SetError(&NuxeoError{}).Put(apiV1 + "/group/" + name)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) DeleteGroup(ctx context.Context, name string, options *nuxeoRequestOptions) error {
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Delete(apiV1 + "/group/" + name)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete group", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (um *UserManager) SearchGroup(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
	reqQuery := "q=" + url.QueryEscape(query)
	if paginationQuery := paginationOptions.QueryParams(); paginationQuery != "" {
		reqQuery += "&" + paginationQuery
	}

	res, err := um.client.NewRequest(ctx, options).SetResult(&Groups{}).SetError(&NuxeoError{}).Get(apiV1 + "/group/search?" + reqQuery)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to search group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Groups), nil
}

func (um *UserManager) AttachGroupToUser(ctx context.Context, idOrGroupName string, idOrUsername string, options *nuxeoRequestOptions) (*Group, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Post(apiV1 + "/group/" + idOrGroupName + "/user/" + idOrUsername)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to attach group to user", slog.String("error", err.Error()))
		return nil, err
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
	res, err := um.client.NewRequest(ctx, options).SetResult(&Users{}).SetError(&NuxeoError{}).Get(urlPath)
	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group member users", slog.String("error", err.Error()))
		return nil, err
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
	res, err := um.client.NewRequest(ctx, options).SetResult(&Groups{}).SetError(&NuxeoError{}).Get(urlPath)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group member groups", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Groups), nil
}

///////////////
//// USERS ////
///////////////

func (um *UserManager) FetchUser(ctx context.Context, id string, options *nuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&User{}).SetError(&NuxeoError{}).Get(apiV1 + "/user/" + id)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

func (um *UserManager) CreateUser(ctx context.Context, user User, options *nuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx, options).SetBody(user).SetResult(&User{}).SetError(&NuxeoError{}).Post(apiV1 + "/user")

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create user", slog.String("error", err.Error()))
		return nil, err
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
	res, err := um.client.NewRequest(ctx, options).SetBody(updatePayload).SetResult(&User{}).SetError(&NuxeoError{}).Put(apiV1 + "/user/" + idOrUsername)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

func (um *UserManager) DeleteUser(ctx context.Context, idOrUsername string, options *nuxeoRequestOptions) error {
	res, err := um.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Delete(apiV1 + "/user/" + idOrUsername)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete user", slog.String("error", err.Error()))
		return err
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
	res, err := um.client.NewRequest(ctx, options).SetResult(&Users{}).SetError(&NuxeoError{}).Get(apiV1 + "/user/search?" + reqQuery)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to search users", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Users), nil
}

func (um *UserManager) AddUserToGroup(ctx context.Context, idOrUsername string, idOrGroupName string, options *nuxeoRequestOptions) (*User, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&User{}).SetError(&NuxeoError{}).Post(apiV1 + "/user/" + idOrUsername + "/group/" + idOrGroupName)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to add user to group", slog.String("error", err.Error()))
		return nil, err
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
		um.logger.Error("Failed to fetch current user login info", slog.String("error", err.Error()))
		return nil, err
	}

	// then fetch the user details
	return um.FetchUser(ctx, loginInfo.Username, nil)
}

func (um *UserManager) FetchWorkflowInstances(ctx context.Context, options *nuxeoRequestOptions) (*Workflows, error) {
	res, err := um.client.NewRequest(ctx, options).SetResult(&Workflows{}).SetError(&NuxeoError{}).Get(apiV1 + "/workflow")

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch workflow instances", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflows), nil
}

func (um *UserManager) StartWorkflowInstance(ctx context.Context, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	// Ensure entity-type is set
	workflow.entity.EntityType = "workflow"
	res, err := um.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&Workflow{}).SetError(&NuxeoError{}).Post(apiV1 + "/workflow")

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to start workflow instance", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflow), nil
}
