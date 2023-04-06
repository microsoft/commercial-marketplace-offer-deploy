package operations

import "fmt"

type OperationResult string

const (
	OperationResultAccepted OperationResult = "Accepted"
)

type OperationType string

const (
	OperationStartDeployment OperationType = "StartDeployment"
	OperationDryRun          OperationType = "DryRun"
	OperationUnknown         OperationType = "Unknown"
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
