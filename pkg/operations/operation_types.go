package operations

type OperationType string

const (
	StartDeployment  OperationType = "StartDeployment"
	DryRunDeployment OperationType = "DryRun"
)

func (o OperationType) String() *string {
	stringValue := string(o)
	return &stringValue
}
