package blockless

// FunctionManifest describes some important configuration options for a Blockless function.
type FunctionManifest struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Function    Function      `json:"function,omitempty"`
	Deployment  Deployment    `json:"deployment,omitempty"`
	Runtime     Runtime       `json:"runtime,omitempty"`
	Cached      bool          `json:"cached,omitempty"`
	Hooks       []interface{} `json:"hooks,omitempty"`
	FSRootPath  string        `json:"fs_root_path,omitempty"`
	Entry       string        `json:"entry,omitempty"`
	ContentType string        `json:"contentType,omitempty"`
	Permissions []string      `json:"permissions,omitempty"`
}

// Runtime is here to support legacy manifests.
type Runtime struct {
	Checksum string `json:"checksum,omitempty"`
	URL      string `json:"url,omitempty"`
}

// Function represents a Blockless function that can be executed.
type Function struct {
	ID         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Version    string   `json:"version,omitempty"`
	Runtime    string   `json:"runtime,omitempty"`
	Extensions []string `json:"extensions,omitempty"`
}

type Deployment struct {
	CID         string    `json:"cid,omitempty"`
	Checksum    string    `json:"checksum,omitempty"`
	URI         string    `json:"uri,omitempty"`
	Methods     []Methods `json:"methods,omitempty"`
	Aggregation string    `json:"aggregation,omitempty"`
	Nodes       int       `json:"nodes,omitempty"`
	File        string    `json:"file,omitempty"`
}

type Methods struct {
	Name       string      `json:"name,omitempty"`
	Entry      string      `json:"entry,omitempty"`
	Arguments  []Parameter `json:"arguments,omitempty"`
	EnvVars    []Parameter `json:"envvars,omitempty"`
	ResultType string      `json:"result_type,omitempty"`
}

// Parameter represents a generic name-value pair.
type Parameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
