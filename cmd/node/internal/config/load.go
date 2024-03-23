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

func Load(args ...string) (*Config, error) {

	// If arguments are not explicitly specified - use arguments we were started with.
	if len(args) == 0 {
		args = os.Args[1:]
	}

	var configPath string

	flags := newCliFlags()
	flags.fs.StringVar(&configPath, "config", "", "path to a config file")

	// General flags.
	flags.stringFlag(roleCfg, DefaultRole)
	flags.uintFlag(concurrencyCfg, DefaultConcurrency)
	flags.stringSliceFlag(bootNodesCfg, nil)
	flags.stringFlag(workspaceCfg, DefaultWorkspace)
	flags.boolFlag(attributesCfg, false)
	flags.stringFlag(peerDBCfg, DefaultPeerDB)
	flags.stringFlag(functionDBCfg, DefaultFunctionDB)
	flags.stringSliceFlag(topicsCfg, nil)

	// Log.
	flags.stringFlag(logLevelCfg, "info")

	// Connectivity flags.
	flags.stringFlag(addressCfg, DefaultAddress)
	flags.uintFlag(portCfg, DefaultPort)
	flags.stringFlag(privateKeyCfg, "")
	flags.boolFlag(websocketCfg, DefaultUseWebsocket)
	flags.uintFlag(websocketPortCfg, DefaultPort)
	flags.stringFlag(dialbackAddressCfg, DefaultAddress)
	flags.uintFlag(dialbackPortCfg, DefaultPort)
	flags.uintFlag(websocketDialbackPortCfg, DefaultPort)

	// Worker node flags.
	flags.stringFlag(runtimePathCfg, "")
	flags.stringFlag(runtimeCLICfg, blockless.RuntimeCLI())
	flags.float64Flag(cpuLimitCfg, 1)
	flags.int64Flag(memLimitCfg, 0)

	// Head node flags.
	flags.stringFlag(restAPICfg, "")

	flags.fs.Parse(args)

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
	// analogous values.
	translate := flagTranslate(flags.groups(), flags.fs, delimiter)

	err := konfig.Load(posflag.ProviderWithFlag(flags.fs, delimiter, konfig, translate), nil)
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

func flagTranslate(flagGroups map[string]configGroup, fs *pflag.FlagSet, delimiter string) func(*pflag.Flag) (string, any) {

	return func(flag *pflag.Flag) (string, any) {
		key := flag.Name
		val := posflag.FlagVal(fs, flag)

		// Should not happen.
		group, ok := flagGroups[key]
		if !ok {
			return key, val
		}

		name := group.Name()
		if name == "" {
			return key, val
		}

		// Log level is a special case because the CLI flag is already prefixed (--log-level).
		if key == logLevelCfg.flag {
			skey := "log" + delimiter + "level"
			return skey, val
		}

		skey := name + delimiter + key
		return skey, val
	}
}
