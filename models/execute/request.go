package execute

// Request describes an execution request.
type Request struct {
	FunctionID string      `json:"function_id"`
	Method     string      `json:"method"`
	Parameters []Parameter `json:"parameters,omitempty"`
	Config     Config      `json:"config"`
}

// Parameter represents an execution parameter, modeled as a key-value pair.
type Parameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Config represents the configurable options for an execution request.
type Config struct {
	Runtime           RuntimeConfig     `json:"runtime,omitempty"`
	Environment       []EnvVar          `json:"env_vars,omitempty"`
	Stdin             *string           `json:"stdin,omitempty"`
	Permissions       []string          `json:"permissions,omitempty"`
	ResultAggregation ResultAggregation `json:"result_aggregation,omitempty"`

	// NodeCount specifies how many nodes should execute this request.
	NodeCount int `json:"number_of_nodes,omitempty"`

	// Threshold (percentage) defines how many nodes should respond with a result to consider this execution successful.
	Threshold float64 `json:"threshold,omitempty"`
}

// EnvVar represents the name and value of the environment variables set for the execution.
type EnvVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ResultAggregation struct {
	Enable     bool        `json:"enable,omitempty"`
	Type       string      `json:"type,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty"`
}
