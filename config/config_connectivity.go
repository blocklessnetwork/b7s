package config

// Connectivity describes the libp2p host that the node will use.
type Connectivity struct {
	Address                 string `koanf:"address"                   flag:"address,a"`
	Port                    uint   `koanf:"port"                      flag:"port,p"`
	PrivateKey              string `koanf:"private-key"               flag:"private-key"`
	DialbackAddress         string `koanf:"dialback-address"          flag:"dialback-address"`
	DialbackPort            uint   `koanf:"dialback-port"             flag:"dialback-port"`
	Websocket               bool   `koanf:"websocket"                 flag:"websocket,w"`
	WebsocketPort           uint   `koanf:"websocket-port"            flag:"websocket-port"`
	WebsocketDialbackPort   uint   `koanf:"websocket-dialback-port"   flag:"websocket-dialback-port"`
	NoDialbackPeers         bool   `koanf:"no-dialback-peers"         flag:"no-dialback-peers"`
	MustReachBootNodes      bool   `koanf:"must-reach-boot-nodes"     flag:"must-reach-boot-nodes"`
	DisableConnectionLimits bool   `koanf:"disable-connection-limits" flag:"disable-connection-limits"`
	ConnectionCount         uint   `koanf:"connection-count"          flag:"connection-count"`
}
