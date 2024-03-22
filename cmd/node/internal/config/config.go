package config

// Config describes the Blockless configuration options.
type Config struct {
	Role           string   `koanf:"role"`
	Concurrency    uint     `koanf:"concurrency"`
	BootNodes      []string `koanf:"boot-nodes"`
	Workspace      string   `koanf:"workspace"`  // TODO: Check - does a head node ever use a workspace?
	LoadAttributes bool     `koanf:"attributes"` // TODO: Head node probably doesn't need attributes..?
	Topics         []string `koanf:"topics"`

	PeerDatabasePath     string `koanf:"peer-db"`
	FunctionDatabasePath string `koanf:"function-db"` // TODO: Head node doesn't need a function database.

	Log          Log          `koanf:"log"`
	Connectivity Connectivity `koanf:"connectivity"`
	Head         Head         `koanf:"head"`
	Worker       Worker       `koanf:"worker"`
}

// Log describes the logging configuration.
type Log struct {
	Level string `koanf:"level"`
}

// Connectivity describes the libp2p host that the node will use.
type Connectivity struct {
	Address               string `koanf:"address"`
	Port                  uint   `koanf:"port"`
	PrivateKey            string `koanf:"private-key"`
	DialbackAddress       string `koanf:"dialback-address"`
	DialbackPort          uint   `koanf:"dialback-port"`
	Websocket             bool   `koanf:"websocket"`
	WebsocketPort         uint   `koanf:"websocket-port"`
	WebsocketDialbackPort uint   `koanf:"websocket-dialback-port"`
}

type Head struct {
	API string `koanf:"rest-api"`
}

type Worker struct {
	RuntimePath        string  `koanf:"runtime-path"`
	RuntimeCLI         string  `koanf:"runtime-cli"`
	CPUPercentageLimit float64 `koanf:"cpu-percentage-limit"`
	MemoryLimitKB      int64   `koanf:"memory-limit"`
}
