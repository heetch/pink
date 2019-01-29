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
	fmt.Fprintf(f, `{"invoker": "binary", "path": "some-path"}`)
	f.Close()
	defer os.Remove(f.Name())

	m, err := LoadManifest(f.Name())
	require.NoError(t, err)
	require.Equal(t, "binary", m.Invoker)
	require.Equal(t, "some-path", m.Path)
}
