package cmdutil

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/spf13/cobra"
)

type RunParameters struct {
	ID              string
	DefinitionFiles []string
	VarsID          string
	EnvID           string
	SkipResultWait  bool
	JUnitOuptutFile string
	RequiredGates   []string
}

func (p RunParameters) Validate(cmd *cobra.Command, args []string) []error {
	errs := []error{}

	hasDefinitionFilesSpecified := p.DefinitionFiles != nil && len(p.DefinitionFiles) > 0
	hasFileIDsSpecified := p.ID != "" && len(p.ID) > 0

	if !hasDefinitionFilesSpecified && !hasFileIDsSpecified {
		errs = append(errs, ParamError{
			Parameter: "resource",
			Message:   "you must specify at least one definition file or resource ID",
		})
	}

	if hasDefinitionFilesSpecified && hasFileIDsSpecified {
		errs = append(errs, ParamError{
			Parameter: "resource",
			Message:   "you cannot specify both a definition file and resource ID",
		})
	}

	if p.JUnitOuptutFile != "" && p.SkipResultWait {
		errs = append(errs, ParamError{
			Parameter: "junit",
			Message:   "--junit option is incompatible with --skip-result-wait option",
		})
	}

	for _, rg := range p.RequiredGates {
		_, err := openapi.NewSupportedGatesFromValue(rg)
		if err != nil {
			errs = append(errs, ParamError{
				Parameter: "required-gates",
				Message:   fmt.Sprintf("invalid option '%s'. %s", rg, validRequiredGatesMsg()),
			})
		}
	}

	return errs
}

func validRequiredGatesMsg() string {
	opts := make([]string, 0, len(openapi.AllowedSupportedGatesEnumValues))
	for _, v := range openapi.AllowedSupportedGatesEnumValues {
		opts = append(opts, string(v))
	}

	return "valid options: " + strings.Join(opts, ", ")
}
