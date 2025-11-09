package main

import (
	"os"

	"github.com/firecrown-media/stax/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
