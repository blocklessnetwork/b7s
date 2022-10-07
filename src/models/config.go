package models

import (
	"strings"

	"github.com/multiformats/go-multiaddr"
)

type AddrList []multiaddr.Multiaddr

func (al *AddrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *AddrList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	addrs := make([]string, 0)
	err := unmarshal(&addrs)
	if err != nil {
		return err
	}

	*al = make([]multiaddr.Multiaddr, 0)
	for _, addr := range addrs {
		a, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return err
		}

		*al = append(*al, a)
	}

	return nil
}

type Config struct {
	Node       ConfigNode       `yaml:"node"`
	Rest       ConfigRest       `yaml:"rest"`
	Protocol   ConfigProtocol   `yaml:"protocol"`
	Logging    ConfigLogging    `yaml:"logging"`
	Repository ConfigRepository `yaml:"repository"`
	Chain      ConfigChain      `yaml:"chain"`
}
type ConfigNode struct {
	Name          string      `yaml:"name"`
	IP            string      `yaml:"ip"`
	Port          string      `yaml:"port"`
	BootNodes     interface{} `yaml:"boot_nodes"`
	UseStaticKeys bool        `yaml:"use_static_keys"`
	WorkspaceRoot string      `yaml:"workspace_root"`
	RuntimePath   string      `yaml:"runtime_path"`
}
type ConfigRest struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}
type ConfigProtocol struct {
	Role string `yaml:"role"`
}
type ConfigLogging struct {
	FilePath string `yaml:"file_path"`
	Level    string `yaml:"level"`
}
type ConfigRepository struct {
	URL string `yaml:"url"`
}
type ConfigChain struct {
	AddressKey string `yaml:"address_key"`
	RPC        string `yaml:"rpc"`
	Home       string `yaml:"home"`
}
