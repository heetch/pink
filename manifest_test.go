package pink

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadManifest(t *testing.T) {
	f, err := ioutil.TempFile("", "manifest")
	require.NoError(t, err)
	fmt.Fprintf(f, `{"invoker": "executable", "path": "some-path", "command": ["some-path"]}`)
	f.Close()
	defer os.Remove(f.Name())

	m, err := LoadManifest(f.Name())
	require.NoError(t, err)
	require.Equal(t, "executable", m.Invoker)
	require.Equal(t, "some-path", m.Path)
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
		{Manifest{Command: []string{"a", "b"}, Invoker: "executable", Path: "somepath"}, false},
	}

	for _, test := range tests {
		err := validateManifest(&test.m)
		require.Equal(t, test.returnsError, err != nil)
	}
}
