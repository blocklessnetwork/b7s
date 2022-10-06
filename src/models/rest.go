package models

type RequestExecute struct {
	FunctionID string                     `json:"function_id"`
	Method     string                     `json:"method"`
	Parameters []RequestExecuteParameters `json:"parameters"`
	Config     ExecutionRequestConfig     `json:"config"`
}
type RequestExecuteEnvVars struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type RequestExecuteParameters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type RequestExecuteResultAggregation struct {
	Enable     bool                       `json:"enable"`
	Type       string                     `json:"type"`
	Parameters []RequestExecuteParameters `json:"parameters"`
}
type ExecutionRequestConfig struct {
	EnvVars           []RequestExecuteEnvVars         `json:"env_vars"`
	NumberOfNodes     int                             `json:"number_of_nodes"`
	ResultAggregation RequestExecuteResultAggregation `json:"result_aggregation"`
}

type ResponseExecute struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Id     string `json:"id"`
	Result string `json:"result"`
}

type RequestFunctionInstall struct {
	Type string `json:"type"`
	Uri  string `json:"uri"`
}

type ResponseInstall struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Result string `json:"result"`
}

type RequestFunctionResponse struct {
	Id string `json:"id"`
}
