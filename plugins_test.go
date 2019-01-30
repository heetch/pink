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

// findPlugins returns the location of all intalled manifest.json files below a given root.
func TestFindPlugin(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", ".pink")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(tmpDir)
	pluginRoot := GetPluginsDir(tmpDir)
	paths := []string{
		filepath.Join(pluginRoot, "system", "run"),
		filepath.Join(pluginRoot, "system", "stop"),
		filepath.Join(pluginRoot, "foo"),
	}
	err = makeManifestPaths(paths)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assertPlugins := func(t *testing.T, root string, expected []string) {
		t.Run(root, func(t *testing.T) {
			plugins, err := findPlugins(root)
			if err != nil {
				t.Fatalf(err.Error())
			}
			require.ElementsMatch(t, expected, plugins)

		})
	}

	assertPlugins(t, pluginRoot, []string{"system", "foo"})
	assertPlugins(t, filepath.Join(pluginRoot, "system"), []string{"run", "stop"})
	assertPlugins(t, filepath.Join(pluginRoot, "system", "run"), []string{})
	assertPlugins(t, filepath.Join(pluginRoot, "system", "stop"), []string{})
	assertPlugins(t, filepath.Join(pluginRoot, "foo"), []string{})
}
