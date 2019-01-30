package pink

import (
	"os"
	"path/filepath"
)

// getPinkDir returns the default state directory for pink
func getPinkDir() string {
	return filepath.Join(os.Getenv("HOME"), ".pink")
}

// findPlugins returns the paths of all manifest.json files of plugins installed below the given root.
func findPlugins(root string) ([]string, error) {
	plugins := []string{}
	wfn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "manifest.json" {
			plugins = append(plugins, path)
		}
		return nil
	}

	err := filepath.Walk(root, wfn)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}
