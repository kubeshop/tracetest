package yaml

import "fmt"

type Transaction struct {
	ID          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Description string            `mapstructure:"description"`
	Env         map[string]string `mapstructure:"env"`
	Steps       []string          `mapstructure:"steps"`
}

func (t Transaction) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("transaction name cannot be empty")
	}

	return nil
}
