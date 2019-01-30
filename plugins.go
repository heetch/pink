package pink

import (
	"io/ioutil"
	"path/filepath"
)

// GetPluginsDir returns the default root of the plugins tree below a given root
func GetPluginsDir(root string) string {
	return filepath.Join(root, "plugins")
}

// findPlugins returns the paths of all manifest.json files of plugins installed below the given root.
func findPlugins(root string) ([]string, error) {
	plugins := []string{}
	// wfn := func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !info.IsDir() && info.Name() == "manifest.json" {
	// 		plugins = append(plugins, path)
	// 	}
	// 	return nil
	// }

	// err := filepath.Walk(root, wfn)
	// if err != nil {
	// 	return nil, err
	// }
	nodes, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		if node.IsDir() {
			plugins = append(plugins, node.Name())
		}
	}
	return plugins, nil
}
