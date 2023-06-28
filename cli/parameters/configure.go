package parameters

import (
	"net/url"

	"github.com/spf13/cobra"
)

type ConfigureParams struct {
	Endpoint string
	Global   bool
}

func (p *ConfigureParams) Validate(cmd *cobra.Command, args []string) []error {
	var errors []error

	if cmd.Flags().Lookup("endpoint").Changed {
		if p.Endpoint == "" {
			errors = append(errors, paramError{
				Parameter: "endpoint",
				Message:   "endpoint cannot be empty",
			})
		} else {
			_, err := url.Parse(p.Endpoint)
			if err != nil {
				errors = append(errors, paramError{
					Parameter: "endpoint",
					Message:   "endpoint is not a valid URL",
				})
			}
		}
	}

	return errors
}
