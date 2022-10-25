package main

import (
	"os"

	"github.com/blocklessnetworking/b7s/src/daemon"
	"github.com/blocklessnetworking/b7s/src/keygen"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setLogging(logType string) {
	if logType == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else if logType == "text" {
		log.SetFormatter(&log.TextFormatter{})
	}

}

func main() {

	// flag values
	var configPath string
	var logType string

	// set the daemon loop to main command
	var rootCmd = cobra.Command{
		Use:   "b7s",
		Short: "Blockless is a peer-to-peer network for compute.",
		Long:  `Blockless is a peer-to-peer network that allows you to earn by sharing your compute power.`,
		Run: func(cmd *cobra.Command, args []string) {
			setLogging(logType)

			log.Info("starting b7s")
			daemon.Run(cmd, args, configPath)
		},
	}

	// generate identity keys
	var keyGenCmd = &cobra.Command{
		Use:   "keygen",
		Short: "Generate a new keypair",
		Long:  `Generate a new keypair`,
		Run: func(cmd *cobra.Command, args []string) {
			setLogging(logType)

			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}

			err = keygen.GenerateKeys(path)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	// add flags
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "path of the config file")
	rootCmd.Flags().StringVarP(&logType, "out", "o", "rich", "output format of logs json, text, rich")

	// add subcommands
	rootCmd.AddCommand(keyGenCmd)

	// execute the cobra loop
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
