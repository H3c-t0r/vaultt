package volume // import "github.com/docker/docker/api/types/volume"

// ----------------------------------------------------------------------------
// Code generated by `swagger generate operation`. DO NOT EDIT.
//
// See hack/generate-swagger-api.sh
// ----------------------------------------------------------------------------

// VolumeCreateBody Volume configuration
// swagger:model VolumeCreateBody
type VolumeCreateBody struct {

	// Name of the volume driver to use.
	// Required: true
	Driver string `json:"Driver"`

	// A mapping of driver options and values. These options are passed directly to the driver and are driver specific.
	// Required: true
	DriverOpts map[string]string `json:"DriverOpts"`

	// User-defined key/value metadata.
	// Required: true
	Labels map[string]string `json:"Labels"`

	// The new volume's name. If not specified, Docker generates a name.
	// Required: true
	Name string `json:"Name"`
}
