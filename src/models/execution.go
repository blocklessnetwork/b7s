package models

type ExecutionRequest struct {
	FunctionID string                       `json:"function_id"`
	Method     string                       `json:"method"`
	Parameters []ExecutionRequestParameters `json:"parameters"`
	Config     ExecutionRequestConfig       `json:"config"`
}
type ExecutionRequestEnvVars struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ExecutionRequestParameters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ExecutionRequestResultAggregation struct {
	Enable     bool                         `json:"enable"`
	Type       string                       `json:"type"`
	Parameters []ExecutionRequestParameters `json:"parameters"`
}
type ExecutionRequestConfig struct {
	EnvVars           []ExecutionRequestEnvVars         `json:"env_vars"`
	NumberOfNodes     int                               `json:"number_of_nodes"`
	ResultAggregation ExecutionRequestResultAggregation `json:"result_aggregation"`
}
