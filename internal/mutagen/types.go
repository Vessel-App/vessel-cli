package mutagen

// SyncSession represents the JSON sync session we get when calling `mutagen sync list`
// This struct exists at github.com/mutagen-io/mutagen/pkg/api/models/synchronization/session.go
// but unmarshalling it to that object results in an error
type SyncSession struct {
	// Identifier is the unique session identifier.
	Identifier string `json:"identifier"`
	// Version is the session version.
	Version uint32 `json:"version"`
	// CreationTime is the session creation timestamp.
	CreationTime string `json:"creationTime"`
	// CreatingVersion is the version of Mutagen that created the session.
	CreatingVersion string `json:"creatingVersion"`
	// Alpha stores the alpha endpoint's configuration and state.
	Alpha Endpoint `json:"alpha"`
	// Beta stores the beta endpoint's configuration and state.
	Beta Endpoint `json:"beta"`
	// Name is the session name.
	Name string `json:"name,omitempty"`
	// Label are the session labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Paused indicates whether or not the session is paused.
	Paused bool `json:"paused"`
}

// ForwardSession represents the JSON sync session we get when calling `mutagen sync list`
// This struct exists at github.com/mutagen-io/mutagen/pkg/api/models/forwarding/session.go
// but unmarshalling it to that object results in an error
type ForwardSession struct {
	// Identifier is the unique session identifier.
	Identifier string `json:"identifier"`
	// Version is the session version.
	Version uint32 `json:"version"`
	// CreationTime is the session creation timestamp.
	CreationTime string `json:"creationTime"`
	// CreatingVersion is the version of Mutagen that created the session.
	CreatingVersion string `json:"creatingVersion"`
	// Source stores the source endpoint's configuration and state.
	Source Endpoint `json:"source"`
	// Destination stores the destination endpoint's configuration and state.
	Destination Endpoint `json:"destination"`
	// Name is the session name.
	Name string `json:"name,omitempty"`
	// Label are the session labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Paused indicates whether or not the session is paused.
	Paused bool `json:"paused"`
}

// Endpoint represents a synchronization endpoint.
// This struct exists at github.com/mutagen-io/mutagen/pkg/api/models/synchronization/session.go
type Endpoint struct {
	// User is the endpoint user.
	User string `json:"user,omitempty"`
	// Host is the endpoint host.
	Host string `json:"host,omitempty"`
	// Port is the endpoint port.
	Port uint16 `json:"port,omitempty"`
	// Path is the synchronization root on the endpoint.
	Path string `json:"path"`
	// Environment is the environment variable map to use for the transport.
	Environment map[string]string `json:"environment,omitempty"`
	// Parameters is the parameter map to use for the transport.
	Parameters map[string]string `json:"parameters,omitempty"`
	// Connected indicates whether or not the controller is currently connected
	// to the endpoint.
	Connected bool `json:"connected"`
}
