package execute

// Request describes an execution request.
type Request struct {
	FunctionID string      `json:"function_id"`
	Method     string      `json:"method"`
	Parameters []Parameter `json:"parameters"`
	Config     Config      `json:"config"`
}

// Parameter represents an execution parameter, modeled as a key-value pair.
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Config represents the configurable options for an execution request.
type Config struct {
	Environment       []EnvVar          `json:"env_vars"`
	NodeCount         int               `json:"number_of_nodes"`
	ResultAggregation ResultAggregation `json:"result_aggregation"`
	Stdin             *string           `json:"stdin"`
	Permissions       []string          `json:"permissions"`
}

// EnvVar represents the name and value of the environment variables set for the execution.
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResultAggregation struct {
	Enable     bool        `json:"enable"`
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}
