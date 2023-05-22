package sdk

type DryRunResponse struct {
	Result DryRunResult `json:"result,omitempty" mapstructure:"result"`
}

type DryRunResult struct {
	Status string        `json:"status,omitempty" mapstructure:"status"`
	Errors []DryRunError `json:"errors,omitempty" mapstructure:"errors"`
}

type DryRunError struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro" mapstructure:"additionalInfo"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro" mapstructure:"code"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro" mapstructure:"message"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro" mapstructure:"target"`

	// READ-ONLY; The error details.
	Details []*DryRunError `json:"details,omitempty" azure:"ro" mapstructure:"details"`
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}
