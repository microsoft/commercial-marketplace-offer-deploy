package deployment

const LookupPrefix = "modm."

type LookupTagKey string

const (
	// reference tag key for events

	// The unique id for modm to identify something
	LookupTagKeyId LookupTagKey = "modm.id"

	// the unique id for an instance of any operation executed by modm against one or more resources
	LookupTagKeyOperationId LookupTagKey = "modm.operation.id"

	// the unique id for an instance of any operation executed by modm against one or more resources
	LookupTagKeyDeploymentId LookupTagKey = "modm.deployment.id"

	// whether or not to send events, if this is not set to true, then the event will not be sent
	LookupTagKeyEvents LookupTagKey = "modm.events"

	// the friendly name of the resource used for logging
	LookupTagKeyName LookupTagKey = "modm.name"

	// the stage id reference. Use is on a resource that's a child of a 1-level parent deployment
	LookupTagKeyStageId LookupTagKey = "modm.stage.id"

	// the stage id reference. Use is on a resource that's a child of a 1-level parent deployment
	LookupTagKeyRetry LookupTagKey = "modm.retry"
)

type LookupTags map[LookupTagKey]*string
