package pink

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
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

// DockerInvoker invokes a plugin as a docker container.
type DockerInvoker struct {
	Client *client.Client
}

// Invoke the docker container based on the given manifest and config.
func (d *DockerInvoker) Invoke(ctx context.Context, m *Manifest, cfg *InvokerConfig) error {
	resp, err := d.Client.ContainerCreate(ctx,
		&container.Config{
			Image:        m.ImageURL,
			Cmd:          cfg.Args,
			Tty:          true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
		}, &container.HostConfig{
			AutoRemove: true,
		}, nil, "")
	if err != nil {
		return errors.Wrapf(err, "unable to create container '%s'", m.ImageURL)
	}

	err = d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to run container '%s'", m.ImageURL)
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rd, err := d.Client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			Follow:     true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", errors.Wrapf(err, "unable to run stream container logs for image '%s'", m.ImageURL))
			return
		}
		defer rd.Close()

		io.Copy(os.Stdout, rd)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rd, err := d.Client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStderr: true,
			Follow:     true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", errors.Wrapf(err, "unable to run stream container logs for image '%s'", m.ImageURL))
			return
		}
		defer rd.Close()

		io.Copy(os.Stderr, rd)
	}()

	_, err = d.Client.ContainerWait(ctx, resp.ID)
	return errors.Wrapf(err, "an error occurs while running container '%s'", m.ImageURL)
}
