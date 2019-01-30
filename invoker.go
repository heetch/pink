package pink

import "context"

// Invoker invokes a plugin and passes it the given configuration.
type Invoker interface {
	Invoke(ctx context.Context, m *Manifest, cfg *InvokerConfig) error
}

// InvokerConfig contains information passed to the invoked plugin.
type InvokerConfig struct {
	Args []string
	Env  []string
}
