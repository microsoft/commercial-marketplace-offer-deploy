package operation

import "fmt"

// OperationType is an enum for the type of operation
type OperationType string

const (
	TypeStartDeployment OperationType = "startDeployment"
	TypeRetryDeployment OperationType = "retryDeployment"
	TypeRetryStage      OperationType = "retryStage"
	TypeDryRun          OperationType = "dryRun"
	TypeUnknown         OperationType = "unknown"
)

func (o OperationType) String() string {
	stringValue := string(o)
	return stringValue
}

func Type(o string) (OperationType, error) {
	// TODO: return tuple with error
	switch o {
	case TypeDryRun.String():
		return TypeDryRun, nil
	case TypeStartDeployment.String():
		return TypeStartDeployment, nil
	case TypeRetryDeployment.String():
		return TypeRetryDeployment, nil
	case TypeRetryStage.String():
		return TypeRetryStage, nil
	default:
		return TypeUnknown, fmt.Errorf("unknown operation type %s", o)
	}
}
