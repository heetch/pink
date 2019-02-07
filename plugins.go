package pink

import (
	"fmt"
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

func isInvokable(path string) (bool, error) {
	nodes, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}
	foundManifest := false
	for _, n := range nodes {
		if !n.IsDir() && n.Name() == "manifest.json" {
			foundManifest = true
			break
		}
	}
	return foundManifest, nil
}

// dispatchCommand finds the manifest for the correct command
func dispatchCommand(args []string, path []string, root string) (string, []string, error) {
	subpath := filepath.Join(path...)
	fullPath := filepath.Join(root, subpath)
	inv, err := isInvokable(fullPath)
	if err != nil {
		return "", nil, err
	}
	if inv {
		fmt.Printf("Invoking %+v with %+v", path, args)
		return filepath.Join(fullPath, "manifest.json"), args, nil
	}
	path = append(path, args[0])
	return dispatchCommand(args[1:], path, root)
}
