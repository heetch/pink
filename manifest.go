package pink

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// A Manifest describes a plugin information.
type Manifest struct {
	Invoker string `json:"invoker"`
	Path    string `json:"path"`
}

// LoadManifest reads the file found at the given path and decodes it into a manifest.
// The file must be JSON formatted.
func LoadManifest(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load manifest at path '%s'", path)
	}
	defer f.Close()

	var m Manifest
	err = json.NewDecoder(f).Decode(&m)
	return &m, errors.Wrap(err, "unable to decode manifest content")
}
