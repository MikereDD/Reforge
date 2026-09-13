package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/MikereDD/Reforge/internal/catalog"
	"github.com/MikereDD/Reforge/internal/resolver"
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
		"\nResolving %s from official upstream infrastructure...\n",
		entry.Name,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		20*time.Second,
	)
	defer cancel()

	registry := resolver.New(nil)

	target, err := registry.Resolve(ctx, entry)
	if err != nil {
		if errors.Is(err, resolver.ErrUnsupported) {
			fmt.Fprintf(
				os.Stdout,
				"\nResolver for %s is not implemented yet.\n",
				entry.Name,
			)
			return
		}

		fmt.Fprintf(os.Stderr, "reforge: %v\n", err)
		os.Exit(1)
	}

	ui.PrintResolvedTarget(os.Stdout, target)
}
