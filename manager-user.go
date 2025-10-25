package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo/internal"
)

type userManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

////////////////
//// GROUPS ////
////////////////

func (um *userManager) FetchGroup(ctx context.Context, name string, options *nuxeoRequestOptions) (*entityGroup, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetResult(&entityGroup{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityGroup), nil
}

func (um *userManager) CreateGroup(ctx context.Context, group entityGroup, options *nuxeoRequestOptions) (*entityGroup, error) {
	path := internal.PathApiV1 + "/group"
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&entityGroup{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityGroup), nil
}

func (um *userManager) UpdateGroup(ctx context.Context, name string, group entityGroup, options *nuxeoRequestOptions) (*entityGroup, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&entityGroup{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityGroup), nil
}

func (um *userManager) DeleteGroup(ctx context.Context, name string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetResult(&entityGroup{}).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete group", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (um *userManager) SearchGroup(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*entityGroups, error) {
	path := internal.PathApiV1 + "/group/search"

	params := url.Values{}
	params.Add("q", query)
	params = internal.MergeUrlValues(params, paginationOptions.QueryParams())
	path += "?" + params.Encode()

	res, err := um.client.NewRequest(ctx, options).SetResult(&entityGroups{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to search group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityGroups), nil
}

func (um *userManager) AttachGroupToUser(ctx context.Context, idOrGroupName string, idOrUsername string, options *nuxeoRequestOptions) (*entityGroup, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(idOrGroupName) + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetResult(&entityGroup{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to attach group to user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityGroup), nil
}

func (um *userManager) FetchGroupMemberUsers(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*entityUsers, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(idOrGroupName) + "/@users"

	if queryPagination := paginationOptions.QueryParams(); queryPagination != nil {
		path += "?" + queryPagination.Encode()
	}

	res, err := um.client.NewRequest(ctx, options).SetResult(&entityUsers{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group member users", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityUsers), nil
}

func (um *userManager) FetchGroupMemberGroups(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*entityGroups, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(idOrGroupName) + "/@groups"

	if queryPagination := paginationOptions.QueryParams(); queryPagination != nil {
		path += "?" + queryPagination.Encode()
	}

	res, err := um.client.NewRequest(ctx, options).SetResult(&entityGroups{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group member groups", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityGroups), nil
}

///////////////
//// USERS ////
///////////////

func (um *userManager) FetchUser(ctx context.Context, id string, options *nuxeoRequestOptions) (*entityUser, error) {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(id)
	res, err := um.client.NewRequest(ctx, options).SetResult(&entityUser{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityUser), nil
}

func (um *userManager) CreateUser(ctx context.Context, user entityUser, options *nuxeoRequestOptions) (*entityUser, error) {
	path := internal.PathApiV1 + "/user"
	res, err := um.client.NewRequest(ctx, options).SetBody(user).SetResult(&entityUser{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityUser), nil
}

func (um *userManager) UpdateUser(ctx context.Context, idOrUsername string, user entityUser, options *nuxeoRequestOptions) (*entityUser, error) {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetBody(user).SetResult(&entityUser{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityUser), nil
}

func (um *userManager) DeleteUser(ctx context.Context, idOrUsername string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete user", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (um *userManager) SearchUsers(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*entityUsers, error) {
	path := internal.PathApiV1 + "/user/search"

	params := url.Values{}
	params.Add("q", query)
	params = internal.MergeUrlValues(params, paginationOptions.QueryParams())

	path += "?" + params.Encode()

	res, err := um.client.NewRequest(ctx, options).SetResult(&entityUsers{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to search users", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityUsers), nil
}

func (um *userManager) AddUserToGroup(ctx context.Context, idOrUsername string, idOrGroupName string, options *nuxeoRequestOptions) (*entityUser, error) {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername) + "/group/" + url.PathEscape(idOrGroupName)
	res, err := um.client.NewRequest(ctx, options).SetResult(&entityUser{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to add user to group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityUser), nil
}

//////////////////////
//// CURRENT USER ////
//////////////////////

func (um *userManager) FetchCurrentUser(ctx context.Context) (*entityUser, error) {
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

func (um *userManager) FetchWorkflowInstances(ctx context.Context, options *nuxeoRequestOptions) (*entityWorkflows, error) {
	path := internal.PathApiV1 + "/workflow"
	res, err := um.client.NewRequest(ctx, options).SetResult(&entityWorkflows{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch workflow instances", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflows), nil
}

func (um *userManager) StartWorkflowInstance(ctx context.Context, workflow entityWorkflow, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/workflow"
	res, err := um.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&entityWorkflow{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to start workflow instance", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}
