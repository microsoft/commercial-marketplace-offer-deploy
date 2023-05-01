package operations

import "fmt"

type OperationType string

const (
	OperationStartDeployment OperationType = "startDeployment"
	OperationRetryDeployment OperationType = "retryDeployment"
	OperationRetryStage      OperationType = "retryStage"
	OperationDryRun          OperationType = "dryRun"
	OperationUnknown         OperationType = "unknown"
)

// Gets the list of operations
func GetOperations() []OperationType {
	return []OperationType{OperationStartDeployment, OperationDryRun}
}

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
	default:
		return OperationUnknown, fmt.Errorf("unknown operation type %s", o)
	}
}
