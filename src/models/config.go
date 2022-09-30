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
	Protocol struct {
		Role           string `yaml:"role"`
		Seed           int64  `yaml:"seed"`
		ProtocolID     string `yaml:"protocol_id"`
		PeerProtocolID string `yaml:"peer_protocol_id"`
	} `yaml:"protocol"`
	Node struct {
		Name               string   `yaml:"name"`
		IpAddress          string   `yaml:"ip"`
		Port               int      `yaml:"port"`
		Rendezvous         string   `yaml:"rendezvous"`
		BootNodes          AddrList `yaml:"boot_nodes"`
		UseStaticKeys      bool     `yaml:"use_static_keys"`
		ConPath            string   `yaml:"conf_path"`
		CoordinatorAddress string   `yaml:"coordinator_address"`
		CoordinatorPort    int      `yaml:"coordinator_port"`
		CoordinatorID      string   `yaml:"coordinator_id"`
		WorkSpaceRoot      string   `yaml:"workspace_root"`
	} `yaml:"node"`
	Rest struct {
		Port    string `yaml:"port"`
		Address string `yaml:"address"`
	} `yaml:"rest"`
	Repository struct {
		Url string `yaml:"url"`
	} `yaml:"repository"`
	Chain struct {
		Disabled   bool   `yaml:"disabled"`
		AddressKey string `yaml:"address_key"`
		RPC        string `yaml:"rpc"`
		Home       string `yaml:"home"`
	} `yaml:"chain"`
	Logging struct {
		FilePath string `yaml:"file_path"`
		Level    string `yaml:"level"`
	} `yaml:"logging"`
}
