package models

type FunctionManifest struct {
	Function    Function      `json:"function,omitempty"`
	Deployment  Deployment    `json:"deployment,omitempty"`
	Runtime     Runtime       `json:"runtime,omitempty"`
	Cached      bool          `json:"cached,omitempty"`
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Hooks       []interface{} `json:"hooks,omitempty"`
	Description string        `json:"description,omitempty"`
	FsRootPath  string        `json:"fs_root_path,omitempty"`
	Entry       string        `json:"entry,omitempty"`
	ContentType string        `json:"contentType,omitempty"`
	Permissions []string      `json:"permissions,omitempty"`
}

// legacy manifest support
type Runtime struct {
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url,omitempty"`
}

type Function struct {
	ID         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Version    string   `json:"version,omitempty"`
	Runtime    string   `json:"runtime,omitempty"`
	Extensions []string `json:"extensions,omitempty"`
}
type Arguments struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type Envvars struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type Methods struct {
	Name       string      `json:"name,omitempty"`
	Entry      string      `json:"entry,omitempty"`
	Arguments  []Arguments `json:"arguments,omitempty"`
	Envvars    []Envvars   `json:"envvars,omitempty"`
	ResultType string      `json:"result_type,omitempty"`
}
type Deployment struct {
	Cid         string    `json:"cid,omitempty"`
	Checksum    string    `json:"checksum,omitempty"`
	Uri         string    `json:"uri,omitempty"`
	Methods     []Methods `json:"methods,omitempty"`
	Aggregation string    `json:"aggregation,omitempty"`
	Nodes       int       `json:"nodes,omitempty"`
	File        string    `json:"file,omitempty"`
}
