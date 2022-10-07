package models

type RequestExecute struct {
	FunctionId string                     `json:"function_id"`
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
	Stdin             *string                         `json:"stdin"`
}

type ResponseExecute struct {
	Code   string `json:"code"`
	Id     string `json:"id"`
	Result string `json:"result"`
}

type RequestFunctionInstall struct {
	Uri   string `json:"uri"`
	Count int    `json:"count"`
}

type ResponseInstall struct {
	Code string `json:"code"`
}

type RequestFunctionResponse struct {
	Id string `json:"id"`
}
