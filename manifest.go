package pink

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// A Manifest describes a plugin information.
type Manifest struct {
	Invoker string       `json:"invoker"`
	Exec    string       `json:"exec"`
	Docker  DockerConfig `json:"docker"`
	Command []string     `json:"command"`
}

// DockerConfig contains options for customizing the docker invoker.
type DockerConfig struct {
	ImageURL string `json:"image-url"`
	TTY      bool   `json:"bool"`
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
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode manifest content")
	}

	return &m, validateManifest(&m)
}

func validateManifest(m *Manifest) error {
	if len(m.Command) == 0 {
		return errors.New("missing 'command' field in manifest file")
	}

	switch m.Invoker {
	case "executable":
		if len(m.Exec) == 0 {
			return errors.Errorf("missing 'exec' field in manifest file for invoker 'executable'")
		}
	case "docker":
		if len(m.Docker.ImageURL) == 0 {
			return errors.Errorf("missing 'docker.image-url' field in manifest file for invoker 'docker'")
		}
	default:
		return errors.Errorf("unsupported invoker '%s', only 'executable' is currently supported", m.Invoker)
	}

	return nil
}
