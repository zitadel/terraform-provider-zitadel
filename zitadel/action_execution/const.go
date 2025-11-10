package action_execution

const (
	TargetsVar       = "targets"
	ExecutionTypeVar = "execution_type"
	NameVar          = "name"

	// Request, Response & Event block variables
	AllVar = "all"

	// Request / Response block
	ServiceVar = "service"
	MethodVar  = "method"

	// Event block
	EventNameVar  = "event"
	EventGroupVar = "group"

	// Function block
	FunctionNameVar = "name"

	executionTypeRequest  = "request"
	executionTypeResponse = "response"
	executionTypeEvent    = "events"
	executionTypeFunction = "function"
)
