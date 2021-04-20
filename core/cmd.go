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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	service *Service
	conf    *Config
)

func CLICommand(content embed.FS, config *Config) *cobra.Command {
	conf = config
	cobra.OnInitialize(func() {
		var err error
		service, err = createService(runner, content)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error occured while initializing server as a service: %s\n", err)
			os.Exit(1)
		}
	})

	cmd := &cobra.Command{
		Use:   "web-serve [commands]",
		Short: "web-serve - A lightweight production build webapp server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner()
		},
	}

	cmd.AddCommand(startCommand(),
		stopCommand(),
		installCommand(),
		removeCommand(),
		statusCommand())

	return cmd
}

func runner() error {
	writer, closer := CombinedWriter()
	defer closer()
	handler, err := routeConfig(service.fs, writer)
	if err != nil {
		return err
	}
	addr := fmt.Sprintf(":%d", conf.Configuration.Server.Port)

	server := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  DurationCast(conf.Configuration.Server.Timeout.Read, time.Second),
		WriteTimeout: DurationCast(conf.Configuration.Server.Timeout.Write, time.Second),
		IdleTimeout:  DurationCast(conf.Configuration.Server.Timeout.IDLE, time.Second),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to run the server: %s\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer signal.Stop(quit)
	sig := <-quit

	fmt.Printf("signal received: %s. exiting... ", sig.String())
	return nil
}

func startCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ServiceState.startService = true
			status, err := service.manageService()
			fmt.Printf("%s\n", status)
			return err
		},
	}
	return cmd
}
func installCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install server as a os service",
		RunE: func(cmd *cobra.Command, args []string) error {
			ServiceState.installService = true
			status, err := service.manageService()
			fmt.Printf("%s\n", status)
			return err
		},
	}
	return cmd
}
func statusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "server current status running/stopped",
		RunE: func(cmd *cobra.Command, args []string) error {
			ServiceState.status = true
			status, err := service.manageService()
			fmt.Printf("%s\n", status)
			return err
		},
	}
	return cmd
}
func stopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ServiceState.stopService = true
			status, err := service.manageService()
			fmt.Printf("%s\n", status)
			return err
		},
	}
	return cmd
}

func removeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "remove/uninstall the service related to the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ServiceState.removeService = true
			status, err := service.manageService()
			fmt.Printf("%s\n", status)
			return err
		},
	}
	return cmd
}
