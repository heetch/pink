package pink

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// GetPluginsDir returns the default root of the plugins tree below a given root.
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

// dispatchCommand finds the manifest for the correct command, or instructs us to run "help".  If neither case is appropriate, an error is returned.
func DispatchCommand(args []string, path []string, root string) (string, []string, bool, error) {
	subpath := filepath.Join(path...)
	fullPath := filepath.Join(root, subpath)
	inv, err := isInvokable(fullPath)
	if err != nil {
		pErr, ok := err.(*os.PathError)
		if ok && pErr.Err == syscall.ENOENT {
			return "", nil, false, fmt.Errorf("No plug-in called %q is installed", strings.Join(path, " "))
		}
		return "", nil, false, err
	}
	if inv {
		return filepath.Join(fullPath, "manifest.json"), args, false, nil
	}
	// We've hit the end of the args without finding a manifest, we should run help!
	if len(args) == 0 || args[0] == "-h" {
		return fullPath, args, true, nil
	}
	path = append(path, args[0])
	return DispatchCommand(args[1:], path, root)
}
