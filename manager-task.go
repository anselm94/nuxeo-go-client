package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// taskManager provides methods to interact with Nuxeo workflow tasks via the REST API.
//
// Supports querying, fetching, reassigning, delegating, and completing tasks.
type taskManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

// FetchTasks queries workflow tasks for a user, workflow instance, or workflow model.
//
// Returns all matching tasks. See https://doc.nuxeo.com/nxdoc/workflow-task-endpoints/#task
func (t *taskManager) FetchTasks(ctx context.Context, userId, workflowInstanceId, workflowModelName string, options *nuxeoRequestOptions) (*Tasks, error) {
	path := internal.PathApiV1 + "/task"
	params := url.Values{}
	if userId != "" {
		params.Add("userId", userId)
	}
	if workflowInstanceId != "" {
		params.Add("workflowInstanceId", workflowInstanceId)
	}
	if workflowModelName != "" {
		params.Add("workflowModelName", workflowModelName)
	}
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	res, err := t.client.NewRequest(ctx, options).SetResult(&Tasks{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to fetch tasks", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Tasks), nil
}

// FetchTask retrieves a workflow task by its ID.
//
// Returns the task if found. See https://doc.nuxeo.com/nxdoc/workflow-task-endpoints/#task
func (t *taskManager) FetchTask(ctx context.Context, taskId string, options *nuxeoRequestOptions) (*Task, error) {
	path := internal.PathApiV1 + "/task/" + url.PathEscape(taskId)
	res, err := t.client.NewRequest(ctx, options).SetResult(&Task{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to fetch task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Task), nil
}

// ReassignTask reassigns a workflow task to new actors, optionally with a comment.
//
// Returns the updated task. See https://doc.nuxeo.com/nxdoc/workflow-task-endpoints/#task
func (t *taskManager) ReassignTask(ctx context.Context, taskId string, actors string, comment string, options *nuxeoRequestOptions) (*Task, error) {
	path := internal.PathApiV1 + "/task/" + url.PathEscape(taskId) + "/reassign"
	params := url.Values{}
	if actors != "" {
		params.Add("actors", actors)
	}
	if comment != "" {
		params.Add("comment", comment)
	}
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	res, err := t.client.NewRequest(ctx, options).SetResult(&Task{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to reassign task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Task), nil
}

// DelegateTask delegates a workflow task to other actors, optionally with a comment.
//
// Returns the updated task. See https://doc.nuxeo.com/nxdoc/workflow-task-endpoints/#task
func (t *taskManager) DelegateTask(ctx context.Context, taskId string, actors string, comment string, options *nuxeoRequestOptions) (*Task, error) {
	path := internal.PathApiV1 + "/task/" + url.PathEscape(taskId) + "/delegate"
	params := url.Values{}
	if actors != "" {
		params.Add("actors", actors)
	}
	if comment != "" {
		params.Add("comment", comment)
	}
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	res, err := t.client.NewRequest(ctx, options).SetResult(&Task{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to delegate task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Task), nil
}

// TaskCompletionRequest represents the payload to complete a workflow task, including variables and comments.
type TaskCompletionRequest struct {
	Id        string           `json:"id"`
	Comment   string           `json:"comment,omitempty"`
	Variables map[string]Field `json:"variables,omitempty"`
}

// CompleteTask completes a workflow task with the specified action and payload.
//
// Returns the updated task. See https://doc.nuxeo.com/nxdoc/workflow-task-endpoints/#task
func (t *taskManager) CompleteTask(ctx context.Context, taskId string, action string, request TaskCompletionRequest, options *nuxeoRequestOptions) (*Task, error) {
	path := internal.PathApiV1 + "/task/" + url.PathEscape(taskId) + "/" + url.PathEscape(action)
	res, err := t.client.NewRequest(ctx, options).SetBody(request).SetResult(&Task{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to complete task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Task), nil
}
