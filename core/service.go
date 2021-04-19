package core

import (
	"embed"
	"github.com/takama/daemon"
	"runtime"
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
		d, err = daemon.New("test-api", "Test API", daemon.GlobalDaemon)
	} else {
		d, err = daemon.New("test-api", "Test API", daemon.SystemDaemon)
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
		return service.Start()
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
