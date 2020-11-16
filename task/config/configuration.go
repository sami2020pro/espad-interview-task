package config 

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Redis struct {
		Host string `json:"host"`
		Port string `json:"port"`
		Password string `json:"password"`
	} `json:"redis"`
	Options struct {
		Schema string `json:"schema"`
		Prefix string `json:"prefix"`
	} `json:"option"`
} 

func FromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config 
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
