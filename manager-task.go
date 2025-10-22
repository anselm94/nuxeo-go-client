package nuxeo

import (
	"context"
	"log/slog"
)

type TaskManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (t *TaskManager) FetchTasks(ctx context.Context) (*Tasks, error) {
	return nil, nil
}

func (t *TaskManager) FetchTask(ctx context.Context, taskId string) (*Task, error) {
	return nil, nil
}

func (t *TaskManager) ReassignTask(ctx context.Context, taskId string, actors string, comment string) (*Task, error) {
	return nil, nil
}

func (t *TaskManager) DelegateTask(ctx context.Context, taskId string, actors string, comment string) (*Task, error) {
	return nil, nil
}

type TaskCompletionRequest struct {
	Id        string         `json:"id"`
	Comment   string         `json:"comment,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`
}

func (t *TaskManager) CompleteTask(ctx context.Context, taskId string, action string, request TaskCompletionRequest) (*Task, error) {
	return nil, nil
}
