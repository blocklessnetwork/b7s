package config

// Config describes the Blockless configuration options.
type Config struct {
	Log          Log     `yaml:"log"`
	DatabasePath string  `yaml:"db-path"`
	Node         Node    `yaml:"node"`
	Workspace    string  `yaml:"workspace"`
	Execute      Execute `yaml:"execute"`
}

// Node describes the configuration options for the Blockless node.
type Node struct {
	Role      string   `yaml:"role"`
	Host      Host     `yaml:"host"`
	API       string   `yaml:"rest-api"`
	BootNodes []string `yaml:"boot-nodes"`
}

// Host describes the libp2p host that the node will use.
type Host struct {
	Address    string `yaml:"address"`
	Port       uint   `yaml:"port"`
	PrivateKey string `yaml:"private-key"`
}

// Log describes the logging configuration.
type Log struct {
	Level string `yaml:"level"`
}

// Execute describes the configuration options for the Blockless worker node.
type Execute struct {
	Runtime string `yaml:"runtime"`
}
