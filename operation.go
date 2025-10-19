package nuxeo

import "context"

// Operation represents a Nuxeo operation.
type Operation struct {
	ID     string
	Params map[string]any
	Input  any
}

// NewOperation creates a new Operation instance.
func NewOperation(id string) *Operation {
	return &Operation{
		ID:     id,
		Params: make(map[string]any),
	}
}

// SetParam sets a parameter for the operation.
func (o *Operation) SetParam(key string, value any) {
	o.Params[key] = value
}

// SetInput sets the input for the operation.
func (o *Operation) SetInput(input any) {
	o.Input = input
}

// Execute runs the operation using the client.
func (o *Operation) Execute(ctx context.Context, client *NuxeoClient) (any, error) {
	// TODO: Implement operation execution logic
	return nil, nil
}
