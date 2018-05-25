//
//  Copyright 2017 Adobe.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

// Contains functions to read/write application configuration.
package config

import (
	"io/ioutil"
	"os"

	"path/filepath"

	"github.com/cloudfoundry-attic/jibber_jabber"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const configPath = ".config/ccd"
const configFile = "config.yml"

type Config struct {
	Scope           string `yaml:"scope"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	ClientId        string `yaml:"client-id"`
	Locale          string `yaml:"locale"`
	AuthEndpoint    string `yaml:"auth-endpoint,omitempty"`
	ServiceEndpoint string `yaml:"service-endpoint,omitempty"`

	path string
}

func configFilePath() string {
	return os.ExpandEnv("${HOME}/" + configPath + "/" + configFile)
}

func configFileExists(path string) (bool, error) {
	switch stat, err := os.Stat(path); {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return true, errors.Wrap(err, "failed to stat config path")
	case !stat.Mode().IsRegular():
		return true, errors.Errorf("path %q not a directory", path)
	default:
		return true, nil
	}
}

// Load will load the configuration file given with the path parameter (a default path will be used if none is given).
// If the configuration does not exist one will be generated with defaults where applicable.
func Load(path string) (*Config, error) {
	cfg := new(Config)

	switch path {
	case "":
		cfg.path = configFilePath()
	default:
		cfg.path = path
	}

	switch exists, err := configFileExists(cfg.path); {
	case err != nil:
		return nil, errors.Wrap(err, "failed to stat config file")
	case !exists:
		cfg.setDefaults()
		return cfg, cfg.Write()
	default:
		fh, err := os.Open(cfg.path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open config file")
		}
		defer fh.Close()

		buf, err := ioutil.ReadAll(fh)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read file")
		}

		return cfg, yaml.Unmarshal(buf, &cfg)
	}
}

// Write will save the configuration.
func (cfg *Config) Write() error {
	if err := os.MkdirAll(filepath.Dir(cfg.path), 0755); err != nil {
		return errors.Wrap(err, "failed to create config path")
	}

	fh, err := os.Create(cfg.path)
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}
	defer fh.Close()

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	_, err = fh.Write(buf)
	return errors.Wrap(err, "failed to write config")
}

func (cfg *Config) setDefaults() {
	cfg.Scope = "openid,creative_cloud"
	cfg.ClientId = ""
	cfg.AuthEndpoint = "https://ims-na1-stg1.adobelogin.com"
	cfg.ServiceEndpoint = "https://scss.adobesc.com/api/v1"
	switch l, err := jibber_jabber.DetectIETF(); {
	case err != nil:
		cfg.Locale = "en_US"
	default:
		cfg.Locale = l
	}
}

func (cfg *Config) Path() string {
	return cfg.path
}
