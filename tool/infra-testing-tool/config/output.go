package config

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
)

func OutputYaml(b []byte) (s string, err error) {
	comp, err := Parse(b)
	if err != nil {
		return s, err
	}

	comp, _ = ToHCL(comp)

	c := GetConfig()

	// marshall to yaml
	b, err = yaml.Marshal(c)
	if err != nil {
		return s, err
	}

	return string(b), nil
}

func OutputJson(b []byte) (s string, err error) {
	comp, err := Parse(b)
	if err != nil {
		return s, err
	}

	comp, _ = ToHCL(comp)

	c := GetConfig()

	// marshall to json (indented)
	b, err = json.MarshalIndent(c, "", "	")
	if err != nil {
		return s, err
	}

	return string(b), nil
}

func OutputHcl(b []byte) (s string, err error) {
	comp, err := Parse(b)
	if err != nil {
		return s, err
	}

	comp, s = ToHCL(comp)
	return s, err
}

func OutputNormal(b []byte) (hcl, y string, err error) {
	comp, err := Parse(b)
	if err != nil {
		return hcl, y, err
	}

	comp, hcl = ToHCL(comp)

	c := GetConfig()

	log.Println("converting config to Yaml")

	b, err = yaml.Marshal(c)
	if err != nil {
		return hcl, y, err
	}

	return hcl, string(b), nil
}
