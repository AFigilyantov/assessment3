package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// add validation
type HTTPServer struct {
	Address      string        `yaml:"address" env-default:":4000"`
	ReadTimeout  time.Duration `yaml:"readtimeout" env-default:"1s"`
	WriteTimeout time.Duration `yaml:"writetimeout" env-default:"1s"`
}

type Config struct {
	HTTPServer HTTPServer `yaml:"http_server"`
	Storage    Storage    `yaml:"storage"`
	JWT        JWT        `yaml:"jwt"`
}

type Storage struct {
	SQLitePath string `yaml:"path" env-default:"db.sql"`
}

type JWT struct {
	Issuer     string        `yaml:"issuer"`
	ExpiresIn  time.Duration `yaml:"expires_in"`
	PublicKey  string        `yaml:"public_key_path"`
	PrivateKey string        `yaml:"private_key_path"`
}

func Parse(p string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(p, c); err != nil {
		return nil, err
	}

	privateKey, err := os.ReadFile(c.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := os.ReadFile(c.JWT.PublicKey)
	if err != nil {
		return nil, err
	}
	c.JWT.PrivateKey = string(privateKey)
	c.JWT.PublicKey = string(publicKey)
	return c, nil
}
