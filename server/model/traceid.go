package model

const TriggerTypeTRACEID TriggerType = "traceid"

type TRACEIDRequest struct {
	ID string `expr_enabled:"true"`
}

type TRACEIDResponse struct {
	ID string
}
