package operations

type OperationType string

const (
	StartDeploymentOperation  OperationType = "StartDeployment"
	DryRunDeploymentOperation OperationType = "DryRun"
)

// Gets the list of operations
func GetOperations() []OperationType {
	return []OperationType{StartDeploymentOperation, DryRunDeploymentOperation}
}

func (o OperationType) String() string {
	stringValue := string(o)
	return stringValue
}

func GetOperationFromString(o string) OperationType {
	// TODO: return tuple with error
	switch o {
	case DryRunDeploymentOperation.String():
		return DryRunDeploymentOperation
	case StartDeploymentOperation.String():
		return StartDeploymentOperation
	default:
		return DryRunDeploymentOperation
	}
}
