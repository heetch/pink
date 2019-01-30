package pink

import (
	"context"
	"os"
	"os/exec"
	"path"
)

// Invoker invokes a plugin and passes it the given configuration.
type Invoker interface {
	Invoke(ctx context.Context, m *Manifest, cfg *InvokerConfig) error
}

// InvokerConfig contains information passed to the invoked plugin.
type InvokerConfig struct {
	Args []string
	Env  []string
}

// ExecutableInvoker invokes a plugin as an executable.
type ExecutableInvoker struct {
	PluginDir string
}

var execCommandContext = exec.CommandContext

// Invoke an executable described by the given manifest. The configuration can be used to
// pass args and environment variables to that executable.
func (e *ExecutableInvoker) Invoke(ctx context.Context, m *Manifest, cfg *InvokerConfig) error {
	execPath := path.Join(e.PluginDir, path.Join(m.Command...), m.Exec)
	cmd := execCommandContext(ctx, execPath, cfg.Args...)
	cmd.Env = append(cmd.Env, cfg.Env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
