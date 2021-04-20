// Package core
/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* 	http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package core

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Info struct {
	Name  string
	Short string
}

type Server struct {
	Port          int
	Mode          string
	Failsafe      bool
	HtmlBlackList []string
	Blacklist     []string
	Page404       string
	Timeout       struct {
		Read  int
		Write int
		IDLE  int
	}
}
type Logging struct {
	Logdir  string
	Logfile string
	Stdout  bool
}

type Config struct {
	Configuration struct {
		Info    Info
		Server  Server
		Logging Logging
	}
}

func (c *Config) Parse(conf string) error {
	viper.SetConfigType("yml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(conf)))

	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("failed to parse configuration YAML: %s\n", err)
	}
	if c.Configuration.Logging.Logdir == "." {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("use path instead of '.' for logdir: %s", err)
		}
		c.Configuration.Logging.Logdir = dir
	}
	return nil
}
