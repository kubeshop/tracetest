package trigger

const TriggerTypeKafka TriggerType = "kafka"

type KafkaRequest struct {
	BrokerURLs      []string             `json:"brokerUrls"`
	Topic           string               `expr_enabled:"true" json:"topic"`
	Headers         []KafkaMessageHeader `json:"headers"`
	Authentication  *KafkaAuthenticator  `json:"authetication,omitempty"`
	MessageKey      string               `expr_enabled:"true" json:"messageKey"`
	MessageValue    string               `expr_enabled:"true" json:"messageValue"`
	SSLVerification bool                 `json:"sslVerification,omitempty"`
}

type KafkaMessageHeader struct {
	Key   string `expr_enabled:"true" json:"key,omitempty"`
	Value string `expr_enabled:"true" json:"value,omitempty"`
}

type KafkaAuthenticator struct {
	Type  string                   `json:"type,omitempty" expr_enabled:"true"`
	Plain *KafkaPlainAuthenticator `json:"plain,omitempty"`
}

type KafkaPlainAuthenticator struct {
	Username string `json:"username,omitempty" expr_enabled:"true"`
	Password string `json:"password,omitempty" expr_enabled:"true"`
}

type KafkaResponse struct {
	Partition string
	Offset    string
	Error     string
}
