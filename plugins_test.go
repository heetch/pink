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

func setUp(t *testing.T) (func(), string) {
	tmpDir, err := ioutil.TempDir("", "pink-test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	cleanUp := func() { os.RemoveAll(tmpDir) }
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
	return cleanUp, pluginRoot

}

// findPlugins returns the location of all intalled manifest.json files below a given root.
func TestFindPlugin(t *testing.T) {
	tearDown, pluginRoot := setUp(t)
	defer tearDown()
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

// IsInvokable indicates the leaves of the command tree
func TestIsInvokable(t *testing.T) {
	tearDown, pluginRoot := setUp(t)
	defer tearDown()
	inv, err := isInvokable(filepath.Join(pluginRoot, "foo"))
	require.NoError(t, err)
	require.True(t, inv)
	inv, err = isInvokable(filepath.Join(pluginRoot, "system"))
	require.False(t, inv)
	inv, err = isInvokable(filepath.Join(pluginRoot, "system", "run"))
	require.True(t, inv)
	inv, err = isInvokable(filepath.Join(pluginRoot, "system", "stop"))
	require.True(t, inv)
	inv, err = isInvokable(filepath.Join(pluginRoot, "system", "run", "up"))
	require.False(t, inv)
}

// DispatchCommand finds the appropriate plugins and passes it the right portion of the arguments
func TestDispatchCommand(t *testing.T) {
	cases := []struct {
		name     string
		inArgs   []string
		err      string
		manifest string
		outArgs  []string
		runHelp  bool
	}{
		{
			name:     "Tier 1 manifest, with zero args",
			inArgs:   []string{"foo"},
			manifest: "foo/manifest.json",
			outArgs:  []string{},
		},
		{
			name:     "Tier 1 manifest, with multiple args",
			inArgs:   []string{"foo", "-v", "-f", "Bernard", "-n", "Cribbins"},
			manifest: "foo/manifest.json",
			outArgs:  []string{"-v", "-f", "Bernard", "-n", "Cribbins"},
		},

		{
			name:     "Tier 2 manifest, with single arg",
			inArgs:   []string{"system", "run", "postgres"},
			manifest: "system/run/manifest.json",
			outArgs:  []string{"postgres"},
		},
		{
			name:   "Non-existant Tier-1 manifest",
			inArgs: []string{"rather", "--name", "Terry Thomas"},
			err:    `No plug-in called "rather" is installed`,
		},
		{
			name:   "Non-existant Tier-2 manifest",
			inArgs: []string{"system", "ding-dong", "--name", "Leslie Phillips"},
			err:    `No plug-in called "system ding-dong" is installed`,
		},
		{
			name:    "Tier-1 with no manifest (run help)",
			inArgs:  []string{"system"},
			runHelp: true,
		},
		{
			name:    "Tier-1, explicitly run help",
			inArgs:  []string{"system", "-h"},
			runHelp: true,
		},
	}

	tearDown, pluginRoot := setUp(t)
	defer tearDown()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			manifest, args, runHelp, err := DispatchCommand(c.inArgs, []string{}, pluginRoot)
			if c.err != "" {
				require.EqualError(t, err, c.err)
				return
			}
			require.NoError(t, err)
			if c.runHelp {
				require.True(t, runHelp)
				return
			}
			require.False(t, runHelp)
			require.Equal(t, c.outArgs, args)
			require.Equal(t, filepath.Join(pluginRoot, c.manifest), manifest)
		})
	}
}
