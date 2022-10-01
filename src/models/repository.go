package models

type FunctionManifest struct {
	Id      string                   `json:"id"`
	Methods []FunctionMethod         `json:"methods"`
	Runtime *FunctionManifestRuntime `json:"runtime"`
}

type FunctionManifestRuntime struct {
	Cid      string `json:"cid"`
	Checksum string `json:"checksum"`
}

type FunctionMethod struct {
	Name       string           `json:"name"`
	Entry      string           `json:"entry"`
	Arguments  []MethodArgument `json:"arguments"`
	ResultType string           `json:"result_type"`
}

type MethodArgument struct {
	Name  string `json:"name"`
	Type_ string `json:"type"`
}
