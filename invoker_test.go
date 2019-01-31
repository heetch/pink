package pink

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

// replace the execCommandContext by a fake one that will run the test binary
// as the executable. The Invoker will run the TestHelperProcess test and expect
// success exit code.
func setFakeExecCommandContext(t *testing.T) func() {
	execCommandContext = func(ctx context.Context, command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.CommandContext(ctx, os.Args[0], cs...)
		cmd.Env = []string{"GO_RUN_HELPER_PROCESS=1"}
		return cmd
	}

	return func() {
		execCommandContext = exec.CommandContext
	}
}

// This test will be run by the Invoker as the targeted executable.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_RUN_HELPER_PROCESS") != "1" {
		return
	}

	// cleanup args, turning
	// /var/folders/xxxx.../pink.test -test.Run=TestHelperProcess -- /..../pink/fake-command --some-flag some-arg
	// into
	// /..../pink/fake-command --some-flag some-arg
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.Equal(t, path.Join(wd, "fake-command"), cmd)
	require.Len(t, args, 2)
	require.Equal(t, "--some-flag", args[0])
	require.Equal(t, "some-arg", args[1])
	require.Equal(t, "C", os.Getenv("A"))
	require.Equal(t, "T", os.Getenv("G"))

	os.Exit(0)
}

func TestExecutableInvoker(t *testing.T) {
	defer setFakeExecCommandContext(t)()

	wd, err := os.Getwd()
	require.NoError(t, err)

	invoker := ExecutableInvoker{
		PluginDir: wd,
	}

	err = invoker.Invoke(
		context.Background(),
		&Manifest{Command: []string{"fake-command"}},
		&InvokerConfig{
			Args: []string{"--some-flag", "some-arg"},
			Env:  []string{"A=C", "G=T"},
		},
	)
	require.NoError(t, err)
}

func TestDockerInvoker(t *testing.T) {
	client, err := client.NewEnvClient()
	require.NoError(t, err)

	invoker := DockerInvoker{
		Client: client,
	}

	err = invoker.Invoke(
		context.Background(),
		&Manifest{ImageURL: "alpine"},
		&InvokerConfig{
			Args: []string{"echo", "hello world"},
		},
	)
	require.NoError(t, err)
}
