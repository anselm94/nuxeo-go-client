package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo/internal"
)

type UserManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

////////////////
//// GROUPS ////
////////////////

func (um *UserManager) FetchGroup(ctx context.Context, name string, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) CreateGroup(ctx context.Context, group Group, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group"
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) UpdateGroup(ctx context.Context, name string, group Group, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) DeleteGroup(ctx context.Context, name string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete group", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (um *UserManager) SearchGroup(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
	path := internal.PathApiV1 + "/group/search"

	params := url.Values{}
	params.Add("q", query)
	params = internal.MergeUrlValues(params, paginationOptions.QueryParams())
	path += "?" + params.Encode()

	res, err := um.client.NewRequest(ctx, options).SetResult(&Groups{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to search group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Groups), nil
}

func (um *UserManager) AttachGroupToUser(ctx context.Context, idOrGroupName string, idOrUsername string, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(idOrGroupName) + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to attach group to user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

func (um *UserManager) FetchGroupMemberUsers(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(idOrGroupName) + "/@users"

	if queryPagination := paginationOptions.QueryParams(); queryPagination != nil {
		path += "?" + queryPagination.Encode()
	}

	res, err := um.client.NewRequest(ctx, options).SetResult(&Users{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group member users", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Users), nil
}

func (um *UserManager) FetchGroupMemberGroups(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(idOrGroupName) + "/@groups"

	if queryPagination := paginationOptions.QueryParams(); queryPagination != nil {
		path += "?" + queryPagination.Encode()
	}

	res, err := um.client.NewRequest(ctx, options).SetResult(&Groups{}).SetError(&NuxeoError{}).Get(path)

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
	path := internal.PathApiV1 + "/user/" + url.PathEscape(id)
	res, err := um.client.NewRequest(ctx, options).SetResult(&User{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

func (um *UserManager) CreateUser(ctx context.Context, user User, options *nuxeoRequestOptions) (*User, error) {
	path := internal.PathApiV1 + "/user"
	res, err := um.client.NewRequest(ctx, options).SetBody(user).SetResult(&User{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

func (um *UserManager) UpdateUser(ctx context.Context, idOrUsername string, user User, options *nuxeoRequestOptions) (*User, error) {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetBody(user).SetResult(&User{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

func (um *UserManager) DeleteUser(ctx context.Context, idOrUsername string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete user", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (um *UserManager) SearchUsers(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
	path := internal.PathApiV1 + "/user/search"

	params := url.Values{}
	params.Add("q", query)
	params = internal.MergeUrlValues(params, paginationOptions.QueryParams())

	path += "?" + params.Encode()

	res, err := um.client.NewRequest(ctx, options).SetResult(&Users{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to search users", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Users), nil
}

func (um *UserManager) AddUserToGroup(ctx context.Context, idOrUsername string, idOrGroupName string, options *nuxeoRequestOptions) (*User, error) {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername) + "/group/" + url.PathEscape(idOrGroupName)
	res, err := um.client.NewRequest(ctx, options).SetResult(&User{}).SetError(&NuxeoError{}).Post(path)

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
	path := internal.PathApiV1 + "/workflow"
	res, err := um.client.NewRequest(ctx, options).SetResult(&Workflows{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch workflow instances", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflows), nil
}

func (um *UserManager) StartWorkflowInstance(ctx context.Context, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	path := internal.PathApiV1 + "/workflow"
	res, err := um.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&Workflow{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to start workflow instance", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflow), nil
}
