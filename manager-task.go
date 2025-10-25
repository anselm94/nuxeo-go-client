package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo/internal"
)

type taskManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (t *taskManager) FetchTasks(ctx context.Context, userId, workflowInstanceId, workflowModelName string, options *nuxeoRequestOptions) (*entityTasks, error) {
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

	res, err := t.client.NewRequest(ctx, options).SetResult(&entityTasks{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to fetch tasks", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityTasks), nil
}

func (t *taskManager) FetchTask(ctx context.Context, taskId string, options *nuxeoRequestOptions) (*entityTask, error) {
	path := internal.PathApiV1 + "/task/" + url.PathEscape(taskId)
	res, err := t.client.NewRequest(ctx, options).SetResult(&entityTask{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to fetch task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityTask), nil
}

func (t *taskManager) ReassignTask(ctx context.Context, taskId string, actors string, comment string, options *nuxeoRequestOptions) (*entityTask, error) {
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

	res, err := t.client.NewRequest(ctx, options).SetResult(&entityTask{}).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to reassign task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityTask), nil
}

func (t *taskManager) DelegateTask(ctx context.Context, taskId string, actors string, comment string, options *nuxeoRequestOptions) (*entityTask, error) {
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

	res, err := t.client.NewRequest(ctx, options).SetResult(&entityTask{}).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to delegate task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityTask), nil
}

type TaskCompletionRequest struct {
	Id        string         `json:"id"`
	Comment   string         `json:"comment,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`
}

func (t *taskManager) CompleteTask(ctx context.Context, taskId string, action string, request TaskCompletionRequest, options *nuxeoRequestOptions) (*entityTask, error) {
	path := internal.PathApiV1 + "/task/" + url.PathEscape(taskId) + "/" + url.PathEscape(action)
	res, err := t.client.NewRequest(ctx, options).SetBody(request).SetResult(&entityTask{}).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		t.logger.Error("Failed to complete task", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityTask), nil
}
