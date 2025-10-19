package internal

// Registry is a simple key-value store for internal use.
type Registry struct {
	store map[string]any
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{store: make(map[string]any)}
}

// Set sets a value in the registry.
func (r *Registry) Set(key string, value any) {
	r.store[key] = value
}

// Get retrieves a value from the registry.
func (r *Registry) Get(key string) any {
	return r.store[key]
}
