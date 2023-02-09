package execute

// Response describes an execution response.
type Response struct {
	Code      string `json:"code"`
	Result    string `json:"result"`
	RequestID string `json:"request_id"`
}
