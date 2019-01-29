package pink

// A Manifest describes a plugin information.
type Manifest struct {
	Invoker string `json:"invoker"`
	Path    string `json:"path"`
}
