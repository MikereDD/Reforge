package main

import (
	"bufio"
	"errors"
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

	reader := bufio.NewReader(os.Stdin)

	ui.PrintInstallerMenu(os.Stdout, c)

	entry, err := ui.SelectInstaller(reader, os.Stdout, c)
	if err != nil {
		if errors.Is(err, ui.ErrCanceled) {
			fmt.Fprintln(os.Stdout, "\nCanceled.")
			return
		}

		fmt.Fprintf(os.Stderr, "reforge: %v\n", err)
		os.Exit(1)
	}

	ui.PrintInstallerDetails(os.Stdout, entry)

	proceed, err := ui.Confirm(reader, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reforge: %v\n", err)
		os.Exit(1)
	}

	if !proceed {
		fmt.Fprintln(os.Stdout, "\nCanceled.")
		return
	}

	fmt.Fprintf(
		os.Stdout,
		"\nSelected %s. Payload resolution is not implemented yet.\n",
		entry.Name,
	)
}
