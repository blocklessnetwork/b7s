package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/config"
	"github.com/blocklessnetworking/b7s/host"
)

const (
	success = 0
	failure = 1
)

const (
	// TODO: Default port for head node is 9527? Move to config if so.
	defaultPort    = 0
	defaultAddress = "0.0.0.0"
)

func main() {
	os.Exit(run())
}

// TODO: Logging format - JSON vs text.
// TODO: Two variants for config loading - look for config file in CWD or explicitely from the flag value.

func run() int {

	var (
		flagAddress  string
		flagConfig   string
		flagLogLevel string
		flagPort     uint

		flagPrivateKey string
	)

	pflag.StringVarP(&flagAddress, "address", "a", defaultAddress, "address to use")
	pflag.StringVarP(&flagConfig, "config", "c", "config.yaml", "path to config file")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level to use")
	pflag.UintVarP(&flagPort, "port", "p", defaultPort, "port number to use - random port if 0")

	pflag.StringVar(&flagPrivateKey, "private-key", "", "private key to use")

	pflag.Parse()

	// Set log level.
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	// Load configuration.
	cfg, err := config.Load(flagConfig)
	if err != nil {
		log.Error().Err(err).Str("config", flagConfig).Msg("could not load configuration")
		return failure
	}
	// TODO: Remove
	_ = cfg

	// Create host.
	host, err := host.New(flagAddress, flagPort, host.WithPrivateKey(flagPrivateKey))
	if err != nil {
		log.Error().Err(err).Str("key", flagPrivateKey).Msg("could not create host")
		return failure
	}

	hostIDs := host.IDs()
	log.Info().Strs("ids", hostIDs).Msg("created host")

	return failure
}
