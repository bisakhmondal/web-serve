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
	"embed"
	"runtime"

	"github.com/takama/daemon"
)

type Service struct {
	daemon.Daemon
	execute func() error
	fs      embed.FS
}

var ServiceState struct {
	startService   bool
	stopService    bool
	installService bool
	removeService  bool
	status         bool
}

func createService(exec func() error, fs embed.FS) (*Service, error) {
	var d daemon.Daemon
	var err error
	if runtime.GOOS == "darwin" {
		d, err = daemon.New(conf.Configuration.Info.Name, conf.Configuration.Info.Short, daemon.GlobalDaemon)
	} else {
		d, err = daemon.New(conf.Configuration.Info.Name, conf.Configuration.Info.Short, daemon.SystemDaemon)
	}
	if err != nil {
		return nil, err
	}
	service := &Service{
		Daemon:  d,
		execute: exec,
		fs:      fs,
	}
	return service, nil
}

func (service *Service) manageService() (string, error) {
	if ServiceState.status {
		return service.Status()
	}
	if ServiceState.installService {
		return service.Install() //"-p", conf.WorkDir)
	}
	if ServiceState.startService {
		iStatus, err := service.Install()
		if err != nil {
			if err != daemon.ErrAlreadyInstalled {
				return iStatus, err
			}
			iStatus = ""
		}
		sStatus, err := service.Start()
		if iStatus != "" {
			sStatus = iStatus + "\n" + sStatus
		}
		return sStatus, err
	} else if ServiceState.stopService {
		return service.Stop()
	}
	if ServiceState.removeService {
		return service.Remove()
	}
	err := service.execute()
	if err != nil {
		return "Unable to start test-api", err
	}
	return "The Test API server exited", nil
}
