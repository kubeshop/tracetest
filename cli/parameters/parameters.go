package parameters

import "fmt"

type paramError struct {
	Parameter string
	Message   string
}

func (pe paramError) Error() string {
	return fmt.Sprintf(`[%s] %s`, pe.Parameter, pe.Message)
}
