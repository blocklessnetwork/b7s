package models

type ExecutorResponse struct {
	Code      string `json:"code"`
	Result    string `json:"result"`
	RequestId string `json:"request_id"`
}
