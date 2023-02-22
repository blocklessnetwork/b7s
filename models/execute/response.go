package execute

// Result describes an execution result.
type Result struct {
	Code      string `json:"code"`
	Result    string `json:"result"`
	RequestID string `json:"request_id"`
}
