package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// UserManager provides methods for managing Nuxeo users and groups via the REST API.
// It supports CRUD operations, search, group membership management, and workflow operations.
type userManager struct {
	client           *NuxeoClient
	logger           *slog.Logger
	operationManager *operationManager
}

////////////////
//// GROUPS ////
////////////////

// FetchGroup retrieves a group by name or ID using the Nuxeo REST API.
// Maps to GET /group/{idOrGroupname}.
func (um *userManager) FetchGroup(ctx context.Context, name string, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

// CreateGroup creates a new group with the given entityGroup payload.
// Maps to POST /group.
func (um *userManager) CreateGroup(ctx context.Context, group Group, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group"
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

// UpdateGroup updates an existing group by name or ID.
// Maps to PUT /group/{idOrGroupname}.
func (um *userManager) UpdateGroup(ctx context.Context, name string, group Group, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetBody(group).SetResult(&Group{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update group", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

// DeleteGroup deletes a group by name or ID.
// Maps to DELETE /group/{idOrGroupname}.
func (um *userManager) DeleteGroup(ctx context.Context, name string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(name)
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete group", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// SearchGroup searches for groups matching the given query and pagination options.
// Maps to GET /group/search?q=...&currentPageIndex=...&pageSize=...
func (um *userManager) SearchGroup(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
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

// AttachGroupToUser adds a user to a group by their respective IDs or names.
// Maps to POST /group/{idOrGroupname}/user/{idOrUsername}.
func (um *userManager) AttachGroupToUser(ctx context.Context, idOrGroupName string, idOrUsername string, options *nuxeoRequestOptions) (*Group, error) {
	path := internal.PathApiV1 + "/group/" + url.PathEscape(idOrGroupName) + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetResult(&Group{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to attach group to user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Group), nil
}

// FetchGroupMemberUsers retrieves users who are members of the specified group.
// Maps to GET /group/{idOrGroupname}/@users.
func (um *userManager) FetchGroupMemberUsers(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
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

// FetchGroupMemberGroups retrieves groups that are members of the specified group.
// Maps to GET /group/{idOrGroupname}/@groups.
func (um *userManager) FetchGroupMemberGroups(ctx context.Context, idOrGroupName string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Groups, error) {
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

// FetchUser retrieves a user by ID or username using the Nuxeo REST API.
// Maps to GET /user/{idOrUsername}.
func (um *userManager) FetchUser(ctx context.Context, id string, options *nuxeoRequestOptions) (*User, error) {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(id)
	res, err := um.client.NewRequest(ctx, options).SetResult(&User{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

// CreateUser creates a new user with the given EntityUser payload.
// Maps to POST /user.
func (um *userManager) CreateUser(ctx context.Context, user User, options *nuxeoRequestOptions) (*User, error) {
	path := internal.PathApiV1 + "/user"
	res, err := um.client.NewRequest(ctx, options).SetBody(user).SetResult(&User{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to create user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

// UpdateUser updates an existing user by ID or username.
// Maps to PUT /user/{idOrUsername}.
func (um *userManager) UpdateUser(ctx context.Context, idOrUsername string, user User, options *nuxeoRequestOptions) (*User, error) {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetBody(user).SetResult(&User{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to update user", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*User), nil
}

// DeleteUser deletes a user by ID or username.
// Maps to DELETE /user/{idOrUsername}.
func (um *userManager) DeleteUser(ctx context.Context, idOrUsername string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/user/" + url.PathEscape(idOrUsername)
	res, err := um.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to delete user", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// SearchUsers searches for users matching the given query and pagination options.
// Maps to GET /user/search?q=...&currentPageIndex=...&pageSize=...
func (um *userManager) SearchUsers(ctx context.Context, query string, paginationOptions *PaginationOptions, options *nuxeoRequestOptions) (*Users, error) {
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

// AddUserToGroup adds a user to a group by their respective IDs or names.
// Maps to POST /user/{idOrUsername}/group/{idOrGroupname}.
func (um *userManager) AddUserToGroup(ctx context.Context, idOrUsername string, idOrGroupName string, options *nuxeoRequestOptions) (*User, error) {
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

// FetchCurrentUser retrieves the currently authenticated user.
// Uses the automation login operation and then fetches user details.
func (um *userManager) FetchCurrentUser(ctx context.Context) (*User, error) {
	// first get the username via the login operation
	loginInfo := &struct {
		Username string `json:"username"`
	}{}

	operationLogin := NewOperation("login")
	if err := um.client.OperationManager().ExecuteInto(ctx, *operationLogin, loginInfo, nil); err != nil {
		um.logger.Error("Failed to fetch current user login info", slog.String("error", err.Error()))
		return nil, err
	}

	// then fetch the user details
	return um.FetchUser(ctx, loginInfo.Username, nil)
}

// FetchWorkflowInstances retrieves workflow instances for the current user.
// Maps to GET /workflow.
func (um *userManager) FetchWorkflowInstances(ctx context.Context, options *nuxeoRequestOptions) (*Workflows, error) {
	path := internal.PathApiV1 + "/workflow"
	res, err := um.client.NewRequest(ctx, options).SetResult(&Workflows{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to fetch workflow instances", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflows), nil
}

// StartWorkflowInstance starts a new workflow instance with the given EntityWorkflow payload.
// Maps to POST /workflow.
func (um *userManager) StartWorkflowInstance(ctx context.Context, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	path := internal.PathApiV1 + "/workflow"
	res, err := um.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&Workflow{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		um.logger.Error("Failed to start workflow instance", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflow), nil
}
