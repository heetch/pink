package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/heetch/pink"
)

// getPinkDir returns the default state directory for pink
func getPinkDir() string {
	return filepath.Join(os.Getenv("HOME"), ".pink")
}

// ensurePinkDirs will create the local filestructure if it doesn't exist
func ensurePinkDirs(root string) error {
	return os.MkdirAll(pink.GetPluginsDir(root), 0700)

}

func main() {
	pd := getPinkDir()
	if err := ensurePinkDirs(pd); err != nil {
		log.Fatalf(err.Error())
	}
	plugD := pink.GetPluginsDir(pd)
	// If we don't pass any arguments default to running help on
	// the root of the plugins directory
	if len(os.Args) == 1 {
		pink.Help(plugD)
		return
	}

	path, args, runHelp, err := pink.DispatchCommand(os.Args[1:], []string{}, plugD)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	if runHelp {
		pink.Help(path)
		return
	}
	err = pink.RunFromManifest(plugD, path, args)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return
}
