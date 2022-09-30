package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	C models.Config
)

func parseYamlFile(file string, o interface{}) error {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if file == "" {
		file = "config.yaml"
	}

	exPath := filepath.Dir(ex)
	configPath := filepath.Join(exPath, file)

	log.WithFields(log.Fields{
		"configPath": configPath,
	}).Info("config path set")

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read file %s failed (%s)", file, err.Error())
	}

	return yaml.Unmarshal(b, o)
}

func Load(cfgFile string) error {
	if err := parseYamlFile(cfgFile, &C); err != nil {
		return err
	}
	return nil
}

func GetPrivKeyFilePath() string {
	return filepath.Join(C.Node.ConPath, "keys/priv.bin")
}

func GetLocalFunctionListFile() string {
	return filepath.Join(C.Node.ConPath, "function_list.json")
}
