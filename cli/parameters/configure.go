package parameters

import (
	"net/url"

	"github.com/spf13/cobra"
)

type ConfigureParams struct {
	Endpoint string
	Global   bool
}

var _ Params = &ConfigureParams{}

func (p *ConfigureParams) Validate(cmd *cobra.Command, args []string) []ParamError {
	var errors []ParamError

	if cmd.Flags().Lookup("endpoint").Changed {
		if p.Endpoint == "" {
			errors = append(errors, ParamError{
				Parameter: "endpoint",
				Message:   "endpoint cannot be empty",
			})
		} else {
			_, err := url.Parse(p.Endpoint)
			if err != nil {
				errors = append(errors, ParamError{
					Parameter: "endpoint",
					Message:   "endpoint is not a valid URL",
				})
			}
		}
	}

	return errors
}
