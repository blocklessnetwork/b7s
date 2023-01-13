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
	if file == "" {
		file = "config.yaml"
	}

	if !filepath.IsAbs(file) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current working directory failed: %v", err)
		}
		file = filepath.Join(cwd, file)
	}

	log.WithFields(log.Fields{
		"configPath": file,
	}).Info("config path set")

	b, err := ioutil.ReadFile(file)
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
