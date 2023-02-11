package api

// RequestExecuteResponse describes the REST API request for the `GetExecuteResponse` call.
type RequestExecuteResponse struct {
	ID string `json:"id"`
}

// RequestFunctionInstall describes the REST API request for the `InstallFunction` call.
type RequestFunctionInstall struct {
	CID   string `json:"cid"`
	URI   string `json:"uri"`
	Count int    `json:"count"`
}
