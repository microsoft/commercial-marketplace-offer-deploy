package operations

import "fmt"

type OperationResult string

const (
	OperationResultAccepted OperationResult = "Accepted"
)

type OperationType string

const (
	StartDeploymentOperation  OperationType = "StartDeployment"
	DryRunDeploymentOperation OperationType = "DryRun"
	UnknownOperation          OperationType = "Unknown"
)

// Gets the list of operations
func GetOperations() []OperationType {
	return []OperationType{StartDeploymentOperation, DryRunDeploymentOperation}
}

func (o OperationType) String() string {
	stringValue := string(o)
	return stringValue
}

func From(o string) (OperationType, error) {
	// TODO: return tuple with error
	switch o {
	case DryRunDeploymentOperation.String():
		return DryRunDeploymentOperation, nil
	case StartDeploymentOperation.String():
		return StartDeploymentOperation, nil
	default:
		return UnknownOperation, fmt.Errorf("unknown operation type %s", o)
	}
}
