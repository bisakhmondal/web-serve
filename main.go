package main

import (
	"embed"
	"fmt"
	"github.com/bisakhmondal/web-serve/core"
	"os"
)

//go:embed html
var html embed.FS

func main() {
	if err := core.CLICommand(html).Execute(); err != nil {
		fmt.Fprintf(os.Stderr,
			"error occured while spinning up the server: %s\n", err)
	}
}
