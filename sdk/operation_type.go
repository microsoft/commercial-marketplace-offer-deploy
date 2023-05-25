package sdk

import "fmt"

// OperationType is an enum for the type of operation
type OperationType string

const (
	OperationDeploy     OperationType = "deploy"
	OperationRetry      OperationType = "retry"
	OperationRetryStage OperationType = "retryStage"
	OperationDryRun     OperationType = "dryRun"
	OperationUnknown    OperationType = "unknown"
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
	case OperationDeploy.String():
		return OperationDeploy, nil
	case OperationRetry.String():
		return OperationRetry, nil
	case OperationRetryStage.String():
		return OperationRetryStage, nil
	default:
		return OperationUnknown, fmt.Errorf("unknown operation type %s", o)
	}
}
