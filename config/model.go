package config

// Config represents the configuration parameters for the node.
type Config struct {
	Node       Node       `yaml:"node"`
	Rest       REST       `yaml:"rest"`
	Protocol   Protocol   `yaml:"protocol"`
	Logging    Logging    `yaml:"logging"`
	Repository Repository `yaml:"repository"`
	Chain      Chain      `yaml:"chain"`
}

type Node struct {
	Name          string   `yaml:"name"`
	BootNodes     []string `yaml:"boot_nodes"`
	WorkspaceRoot string   `yaml:"workspace_root"`
	RuntimePath   string   `yaml:"runtime_path"`
}

type REST struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type Protocol struct {
	Role string `yaml:"role"`
}

type Logging struct {
	FilePath string `yaml:"file_path"`
	Level    string `yaml:"level"`
}

type Repository struct {
	URL string `yaml:"url"`
}

type Chain struct {
	AddressKey string `yaml:"address_key"`
	RPC        string `yaml:"rpc"`
	Home       string `yaml:"home"`
}
