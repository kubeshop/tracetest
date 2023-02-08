package config

import "fmt"

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider")

func (c *Config) DataStore() (*TracingBackendDataStoreConfig, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	selectedStore := c.config.Server.Telemetry.DataStore
	dataStoreConfig, found := c.config.Telemetry.DataStores[selectedStore]

	if selectedStore != "" && !found {
		return nil, ErrInvalidTraceDBProvider
	}

	if !found {
		return nil, nil
	}

	return &dataStoreConfig, nil
}

func (c *Config) IsDataStoreConfigured() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	dataStore, _ := c.DataStore()

	return dataStore != nil
}

func (c *Config) Exporter() (*TelemetryExporterOption, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getExporter(c.config.Server.Telemetry.Exporter)
}

func (c *Config) ApplicationExporter() (*TelemetryExporterOption, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getExporter(c.config.Server.Telemetry.ApplicationExporter)
}

func (c *Config) getExporter(name string) (*TelemetryExporterOption, error) {
	// Exporters are optional: if no name was provided we consider that the user don't want to have them enabled
	if name == "" {
		return nil, nil
	}

	exporterConfig, found := c.config.Telemetry.Exporters[name]
	if !found {
		availableOptions := mapKeys(c.config.Telemetry.DataStores)
		return nil, fmt.Errorf(`invalid exporter option: "%s". Available options: %v`, name, availableOptions)
	}

	return &exporterConfig, nil
}

func mapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
