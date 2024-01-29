package helm

import (
	"encoding/json"
	"fmt"
)

// DefaultNameFormat to use when no nameFormat is supplied
const DefaultNameFormat = `{{ print .kind "_" .metadata.name | snakecase }}`

// VendorOpts are additional properties the consumer of the native func might
// pass.
type VendorOpts struct {
	TemplateOpts

	// CalledFrom is the file that calls helmTemplate. This is used to find the
	// vendored chart relative to this file
	CalledFrom string `json:"calledFrom"`
	// NameTemplate is used to create the keys in the resulting map
	NameFormat string `json:"nameFormat"`
}

func parseOpts(data interface{}) (*VendorOpts, error) {
	c, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// default IncludeCRDs to true, as this is the default in the `helm install`
	// command. Needs to be specified here because the zero value of bool is
	// false.
	opts := VendorOpts{
		TemplateOpts: TemplateOpts{
			IncludeCRDs: true,
		},
	}

	if err := json.Unmarshal(c, &opts); err != nil {
		return nil, err
	}

	// Charts are only allowed at relative paths. Use conf.CalledFrom to find the callers directory
	if opts.CalledFrom == "" {
		return nil, fmt.Errorf("helmTemplate: 'opts.calledFrom' is unset or empty.")
	}

	return &opts, nil
}