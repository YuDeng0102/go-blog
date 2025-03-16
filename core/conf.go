package core

import (
	"log"
	"server/config"
	"server/utils"

	"gopkg.in/yaml.v3"
)

func InitConf() *config.Config {
	c := &config.Config{}
	yamlConf, err := utils.LoadYAML()
	if err != nil {
		log.Fatal("Failed to load configration: %v", err)
	}
	if err = yaml.Unmarshal(yamlConf, c); err != nil {
		log.Fatal("Failed to unmarshal configration: %v", err)
	}
	return c
}
