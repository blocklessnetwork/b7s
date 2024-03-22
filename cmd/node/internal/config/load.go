package config

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func Load() (*Config, error) {

	var configPath string
	fs := pflag.NewFlagSet(flagConfig, pflag.ExitOnError)

	fs.StringVar(&configPath, flagConfig, "", "path to a config file")

	// General node flags.
	fs.StringP(flagRole, "r", DefaultRole, "role this note will have in the Blockless protocol (head or worker)")
	fs.UintP(flagConcurrency, "c", DefaultConcurrency, "maximum number of requests node will process in parallel")
	fs.StringSlice(flagBootNodes, nil, "list of addresses that this node will connect to on startup, in multiaddr format")
	fs.String(flagWorkspace, DefaultWorkspace, "directory that the node can use for file storage")
	fs.Bool(flagAttributes, false, "node should try to load its attribute data from IPFS")
	fs.String(flagPeerDB, DefaultPeerDB, "path to the database used for persisting peer data")
	fs.String(flagFunctionDB, DefaultFunctionDB, "path to the database used for persisting function data")
	fs.StringSlice(flagTopics, nil, "topics node should subscribe to")

	fs.StringP(flagLogLevel, "l", "info", "log level to use")

	// Connectivity flags.
	fs.StringP(flagAddress, "a", DefaultAddress, "address that the b7s host will use")
	fs.UintP(flagPort, "p", DefaultPort, "port that the b7s host will use")
	fs.String(flagPrivateKey, "", "private key that the b7s host will use")
	fs.BoolP(flagWebsocket, "w", DefaultUseWebsocket, "should the node use websocket protocol for communication")
	fs.Uint(flagWebsocketPort, DefaultPort, "port to use for websocket connections")
	fs.StringP(flagDialbackAddress, "", DefaultAddress, "external address that the b7s host will advertise")
	fs.UintP(flagDialbackPort, "", DefaultPort, "external port that the b7s host will advertise")
	fs.UintP(flagWebsocketDialbackPort, "", DefaultPort, "external port that the b7s host will advertise for websocket connections")

	// Head node flags.
	fs.String(flagRestAPI, "", "address where the head node REST API will listen on")

	// Worker node flags.
	fs.String(flagRuntimePath, "", "Blockless Runtime location (used by the worker node)")
	fs.String(flagRuntimeCLI, blockless.RuntimeCLI(), "runtime CLI name (used by the worker node)")
	fs.Float64(flagCPULimit, 1, "amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited")
	fs.Int64(flagMemoryLimit, 0, "memory limit (kB) for Blockless Functions")

	fs.SortFlags = false
	fs.Parse(os.Args[1:])

	delimiter := "."
	konfig := koanf.New(delimiter)

	if configPath != "" {
		err := konfig.Load(file.Provider(configPath), yaml.Parser())
		if err != nil {
			return nil, fmt.Errorf("could not load config file: %w", err)
		}
	}

	// For readability flags have a flat structure - e.g. port or cpu-percentage-limit.
	// For use in config files, we prefer a structured layout, e.g. connectivity=>port or worker=>cpu-percentage-limit.
	// This callback translates the flag names from a flat layout to the structured one, so that koanf knows how to match
	// analogoues values.
	// TODO: This is a bit fragile and assumes a fair amount of responsibility from the dev.
	// We have a tight coupling between flat flag list and the Config structure and the tag values.
	translate := flagTranslate(fs, delimiter)

	err := konfig.Load(posflag.ProviderWithFlag(fs, delimiter, konfig, translate), nil)
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	var cfg Config
	err = konfig.Unmarshal("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal konfig: %w", err)
	}

	return &cfg, nil
}

func flagTranslate(fs *pflag.FlagSet, delimiter string) func(*pflag.Flag) (string, any) {
	return func(flag *pflag.Flag) (string, any) {
		key := flag.Name
		val := posflag.FlagVal(fs, flag)

		switch key {
		// For general flags, we don't have a group prefix.
		default:
			return key, val

		// Connectivity flags:
		case flagAddress,
			flagPort,
			flagPrivateKey,
			flagDialbackAddress,
			flagDialbackPort,
			flagWebsocket,
			flagWebsocketPort,
			flagWebsocketDialbackPort:

			skey := "connectivity" + delimiter + key
			return skey, val

		// Head node flags:
		case flagRestAPI:
			skey := "head" + delimiter + key
			return skey, val

		// Worker node flags:
		case flagRuntimePath, flagRuntimeCLI, flagCPULimit, flagMemoryLimit:
			skey := "worker" + delimiter + key
			return skey, val

		// Log flags:
		case flagLogLevel:
			skey := "log" + delimiter + "level"
			return skey, val
		}
	}
}
