package pink

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeManifestPaths(pluginPaths []string) error {
	for _, pluginPath := range pluginPaths {
		err := os.MkdirAll(pluginPath, 0777)
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(pluginPath, "manifest.json"))
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func TestFindPlugin(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", ".pink")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(tmpDir)
	paths := []string{
		filepath.Join(tmpDir, "plugins", "system", "run"),
		filepath.Join(tmpDir, "plugins", "system", "stop"),
		filepath.Join(tmpDir, "plugins", "foo"),
	}
	err = makeManifestPaths(paths)
	if err != nil {
		t.Fatalf(err.Error())
	}

	plugins, err := findPlugins(tmpDir)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := []string{}
	for _, p := range paths {
		expected = append(expected, filepath.Join(p, "manifest.json"))
	}
	require.ElementsMatch(t, expected, plugins)
}
