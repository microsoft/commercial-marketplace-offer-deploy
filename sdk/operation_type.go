package sdk

import "fmt"

// OperationType is an enum for the type of operation
type OperationType string

const (
	OperationStartDeployment OperationType = "startDeployment"
	OperationRetryDeployment OperationType = "retryDeployment"
	OperationRetryStage      OperationType = "retryStage"
	OperationDryRun          OperationType = "dryRun"
	OperationUnknown         OperationType = "unknown"
)

func (o OperationType) String() string {
	stringValue := string(o)
	return stringValue
}

func Type(o string) (OperationType, error) {
	// TODO: return tuple with error
	switch o {
	case OperationDryRun.String():
		return OperationDryRun, nil
	case OperationStartDeployment.String():
		return OperationStartDeployment, nil
	case OperationRetryDeployment.String():
		return OperationRetryDeployment, nil
	case OperationRetryStage.String():
		return OperationRetryStage, nil
	default:
		return OperationUnknown, fmt.Errorf("unknown operation type %s", o)
	}
}
