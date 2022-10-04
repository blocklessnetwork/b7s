package models

type FunctionManifest struct {
	Id      string                   `json:"id"`
	Methods []FunctionMethod         `json:"methods"`
	Runtime *FunctionManifestRuntime `json:"runtime"`
}

type FunctionManifestRuntime struct {
	Checksum string `json:"checksum"`
	Uri      string `json:"uri"`
}

type FunctionMethod struct {
	Name      string           `json:"name"`
	Entry     string           `json:"entry"`
	Arguments []MethodArgument `json:"arguments"`
}

type MethodArgument struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MethodEnvironment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
