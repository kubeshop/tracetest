package parameters

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type ListParams struct {
	Take          int32
	Skip          int32
	SortBy        string
	SortDirection string
	All           bool
}

func (p ListParams) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.Take < 0 {
		errors = append(errors, paramError{
			Parameter: "take",
			Message:   "take parameter must be greater than 0",
		})
	}

	if p.Skip < 0 {
		errors = append(errors, paramError{
			Parameter: "skip",
			Message:   "skip parameter must be greater than 0",
		})
	}

	if p.SortDirection != "" && p.SortDirection != "asc" && p.SortDirection != "desc" {
		errors = append(errors, paramError{
			Parameter: "sortDirection",
			Message:   "sortDirection parameter must be either asc or desc",
		})
	}

	return errors
}

type ResourceIdParams struct {
	ResourceID string
}

func (p *ResourceIdParams) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.ResourceID == "" {
		errors = append(errors, paramError{
			Parameter: "id",
			Message:   "resource id must be provided",
		})
	}

	return errors
}

type ExportParams struct {
	ResourceIdParams
	OutputFile string
}

func (p *ExportParams) Validate(cmd *cobra.Command, args []string) []error {
	errors := p.ResourceIdParams.Validate(cmd, args)

	if p.OutputFile == "" {
		errors = append(errors, paramError{
			Parameter: "file",
			Message:   "output file must be provided",
		})
	}

	return errors
}

type ApplyParams struct {
	DefinitionFile string
}

func (p *ApplyParams) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.DefinitionFile == "" {
		errors = append(errors, paramError{
			Parameter: "file",
			Message:   "Definition file must be provided",
		})
	}

	return errors
}

type ResourceParams struct {
	ResourceName string
}

var ValidResources = []string{"config", "datastore", "demo", "environment", "pollingprofile", "transaction", "analyzer"}

func (p *ResourceParams) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if len(args) == 0 {
		errors = append(errors, paramError{
			Parameter: "resource",
			Message:   "resource name must be provided",
		})

		return errors
	}

	p.ResourceName = args[0]
	if p.ResourceName == "" {
		errors = append(errors, paramError{
			Parameter: "resource",
			Message:   "resource name must be provided",
		})
	}

	exists := false
	for _, validArg := range ValidResources {
		if validArg == p.ResourceName {
			exists = true
			break
		}
	}

	if !exists {
		errors = append(errors, paramError{
			Parameter: "resource",
			Message:   fmt.Sprintf("resource must be one of %s", strings.Join(ValidResources, ", ")),
		})
	}

	return errors
}
