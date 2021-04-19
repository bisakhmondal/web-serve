package core

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var service *Service

func CLICommand(content embed.FS) *cobra.Command {
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
	router, err := routeConfig(service.fs)
	if err != nil {
		return err
	}
	go router.Run(":8080")

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	sis := <-sig

	fmt.Printf("signal received: %s. exiting... ", sis.String())
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
