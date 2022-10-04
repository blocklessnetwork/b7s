package models

type FunctionManifest struct {
	Function   Function   `json:"function"`
	Deployment Deployment `json:"deployment"`
}
type Function struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	BuildCommand string   `json:"build-command"`
	BuildOutput  string   `json:"build-output"`
	Runtime      string   `json:"runtime"`
	Extensions   []string `json:"extensions"`
}
type Arguments struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Envvars struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Methods struct {
	Name      string      `json:"name"`
	Entry     string      `json:"entry"`
	Arguments []Arguments `json:"arguments"`
	Envvars   []Envvars   `json:"envvars"`
}
type Deployment struct {
	Checksum    string    `json:"checksum"`
	URI         string    `json:"uri"`
	Permission  string    `json:"permission"`
	Methods     []Methods `json:"methods"`
	Aggregation string    `json:"aggregation"`
	Nodes       int       `json:"nodes"`
}
