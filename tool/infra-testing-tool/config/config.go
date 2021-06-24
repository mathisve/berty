package config

import "C"
import (
	"fmt"
	"infratesting/iac"
	"infratesting/iac/components/networking"

	//relay "github.com/libp2p/go-libp2p-circuit"
	"gopkg.in/yaml.v3"
	"log"
)

type Config struct {
	RDVP      []NodeGroup `yaml:"rdvp"`
	Relay     []NodeGroup `yaml:"relay"`
	Bootstrap []NodeGroup `yaml:"bootstrap"`

	Replication []NodeGroup `yaml:"replication"`
	Peer        []NodeGroup `yaml:"peer"`

	Attributes Attributes `yaml:"attributes"`
}

// Attributes are used at infra-compile time to aid with the generation of HCL
type Attributes struct {
	Connections map[string]Connection
	Groups      map[string]Group

	vpc                  networking.Vpc
	connectionComponents map[string][]iac.Component
}

var (
	config = Config{}
)

func init() {
	config.Attributes.Connections = make(map[string]Connection)
	config.Attributes.Groups = make(map[string]Group)
	config.Attributes.connectionComponents = make(map[string][]iac.Component)
}

// Validate function validates the config, gives all NodeGroup's their correct type.
func (c *Config) Validate() error {
	for i := range c.RDVP {

		c.RDVP[i].NodeType = NodeTypeRDVP
		_ = c.RDVP[i].validate()

	}

	for i := range c.Relay {

		c.Relay[i].NodeType = NodeTypeRelay
		_ = c.Relay[i].validate()

	}

	for i := range c.Bootstrap {

		c.Bootstrap[i].NodeType = NodeTypeBootstrap
		_ = c.Bootstrap[i].validate()

	}

	for i := range c.Replication {

		c.Replication[i].NodeType = NodeTypeReplication
		_ = c.Replication[i].validate()

	}

	for i := range c.Peer {

		c.Peer[i].NodeType = NodeTypePeer
		_ = c.Peer[i].validate()

	}
	// TODO
	// add more checks to validate if network topology is correct, etc

	return nil
}

func NewConfig() Config {
	var c = Config{}

	s, err := yaml.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(s))

	return c
}

func GetConfig() Config {
	return config
}
