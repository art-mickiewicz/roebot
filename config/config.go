package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var Config = ConfigSpec{}

type ConfigSpec struct {
	Token  string
	Admins []string
}

func init() {
	bytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(bytes, &Config); err != nil {
		log.Fatalf("Cannot parse config file: %v", err)
	}
}
