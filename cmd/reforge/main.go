package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MikereDD/Reforge/internal/catalog"
	"github.com/MikereDD/Reforge/internal/ui"
)

func main() {
	catalogPath := flag.String(
		"catalog",
		"manifests/catalog.example.json",
		"path to the Reforge installer catalog",
	)

	flag.Parse()

	c, err := catalog.Load(*catalogPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reforge: %v\n", err)
		os.Exit(1)
	}

	ui.PrintInstallerMenu(os.Stdout, c)
}
