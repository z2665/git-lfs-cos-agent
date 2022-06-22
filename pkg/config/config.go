package config

import (
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	SecretID   string `toml:"secretID"`
	SecretKey  string `toml:"secretKey"`
	BucketName string `toml:"bucketName"`
	Region     string `toml:"region"`
}

func LoadConfig(path *string) (Config, error) {
	var conf Config
	if path != nil {
		f, err := os.Open(*path)
		if err != nil {
			return Config{}, err
		}
		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return Config{}, err
		}
		err = toml.Unmarshal(fb, &conf)
		if err != nil {
			return Config{}, err
		}
	}
	return conf, nil

}
