package models

type RequestExecute struct {
	Type string `json:"type"`
	Id   string `json:"id"`
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
