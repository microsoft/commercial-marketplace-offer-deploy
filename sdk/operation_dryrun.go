package sdk

type DryRunResponse struct {
	DryRunResult
}

type DryRunResult struct {
	Status *string `json:"status,omitempty" azure:"ro"`
	//Message *string `json:"message,omitempty" azure:"ro"`
	//Target *string `json:"target,omitempty" azure:"ro"`
	Error *DryRunErrorResponse
}

type DryRunErrorResponse struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro" mapstructure:"additionalInfo"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro" mapstructure:"code"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro" mapstructure:"message"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro" mapstructure:"target"`

	// READ-ONLY; The error details.
	Details []*DryRunErrorResponse `json:"details,omitempty" azure:"ro" mapstructure:"details"`
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}
