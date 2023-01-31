package yaml

import "fmt"

type TRACEIDRequest struct {
	ID string `yaml:"id"`
}

func (t TRACEIDRequest) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	return nil
}
