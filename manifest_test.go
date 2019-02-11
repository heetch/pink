package pink

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadManifest(t *testing.T) {
	t.Run("Executable", func(t *testing.T) {
		f, err := ioutil.TempFile("", "manifest")
		require.NoError(t, err)
		fmt.Fprintf(f, `{"invoker": "executable", "exec": "some-path", "command": ["some-path"]}`)
		f.Close()
		defer os.Remove(f.Name())

		m, err := LoadManifest(f.Name())
		require.NoError(t, err)
		require.Equal(t, "executable", m.Invoker)
		require.Equal(t, "some-path", m.Exec)
	})

	t.Run("Docker", func(t *testing.T) {
		f, err := ioutil.TempFile("", "manifest")
		require.NoError(t, err)
		fmt.Fprintf(f, `{"invoker": "docker", "docker": {"image-url": "some-url", "tty": true}, "command": ["some-path"]}`)
		f.Close()
		defer os.Remove(f.Name())

		m, err := LoadManifest(f.Name())
		require.NoError(t, err)
		require.Equal(t, "docker", m.Invoker)
		require.Equal(t, "some-url", m.Docker.ImageURL)
	})
}

func TestValidateManifest(t *testing.T) {
	tests := []struct {
		m            Manifest
		returnsError bool
	}{
		{Manifest{}, true},
		{Manifest{Invoker: "docker"}, true},
		{Manifest{Command: []string{"a", "b"}, Invoker: "docker"}, true},
		{Manifest{Command: []string{"a", "b"}, Invoker: "executable"}, true},
		{Manifest{Command: []string{"a", "b"}, Invoker: "executable", Exec: "somepath"}, false},
	}

	for _, test := range tests {
		err := validateManifest(&test.m)
		require.Equal(t, test.returnsError, err != nil)
	}
}

// Invoke defers to the correct invoker
func TestInvoke(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	t.Run("executable", func(t *testing.T) {
		execCommandContext = func(ctx context.Context, command string, args ...string) *exec.Cmd {
			require.Equal(t, path.Join(wd, "a/b/somepath"), command)
			return &exec.Cmd{Path: "/bin/echo"}
		}
		defer func() { execCommandContext = exec.CommandContext }()

		m := Manifest{Command: []string{"a", "b"}, Invoker: "executable", Exec: "somepath"}

		require.NoError(t, err)
		err = m.Invoke(context.Background(), wd, nil)
		require.NoError(t, err)

	})

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	tmpPath := filepath.Join(tmpDir, "stdout")
	tmpFile, err := os.Create(tmpPath)
	require.NoError(t, err)
	origStdout := os.Stdout
	os.Stdout = tmpFile
	defer func() { os.Stdout = origStdout }()

	m := Manifest{Invoker: "docker", Docker: DockerConfig{ImageURL: "alpine"}}
	err = m.Invoke(context.Background(), wd, []string{"echo", "hello", "world"})
	require.NoError(t, err)

	tmpFile.Sync()
	tmpFile.Close()

	b, err := ioutil.ReadFile(tmpPath)
	require.NoError(t, err)

	require.Equal(t, "hello world", strings.TrimSpace(string(b)))

}
