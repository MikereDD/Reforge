package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/MikereDD/Reforge/internal/catalog"
	"github.com/MikereDD/Reforge/internal/fetch"
	"github.com/MikereDD/Reforge/internal/resolver"
	"github.com/MikereDD/Reforge/internal/ui"
	"github.com/MikereDD/Reforge/internal/verify"
)

func defaultArchTrustDB() string {
	if runtime.GOOS == "linux" {
		return "/etc/pacman.d/gnupg"
	}

	return ""
}

func defaultArchKeyringSource() string {
	if runtime.GOOS == "linux" {
		return "/usr/share/pacman/keyrings/archlinux.gpg"
	}

	return ""
}

func main() {
	catalogPath := flag.String(
		"catalog",
		"manifests/catalog.example.json",
		"path to the Reforge installer catalog",
	)

	verifyArtifacts := flag.Bool(
		"verify",
		false,
		"download and cryptographically verify resolved boot artifacts",
	)

	pacmanKeyBinary := flag.String(
		"pacman-key",
		"pacman-key",
		"pacman-key binary used for Arch signature verification",
	)

	archTrustDB := flag.String(
		"arch-trust-db",
		defaultArchTrustDB(),
		"path to the initialized Arch pacman GnuPG trust database",
	)

	archKeyringSource := flag.String(
		"arch-keyring-source",
		defaultArchKeyringSource(),
		"path to the packaged Arch Linux keyring source",
	)

	workDir := flag.String(
		"work-dir",
		filepath.Join(
			os.TempDir(),
			"reforge",
			"artifacts",
		),
		"artifact verification workspace",
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

	resolveCtx, cancelResolve := context.WithTimeout(
		context.Background(),
		20*time.Second,
	)
	defer cancelResolve()

	registry := resolver.New(nil)

	target, err := registry.Resolve(resolveCtx, entry)
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

	if !*verifyArtifacts {
		return
	}

	if target.Provider != "archlinux" {
		ui.PrintVerificationPreflightBlocked(
			os.Stdout,
			fmt.Errorf(
				"verification provider %q is not implemented",
				target.Provider,
			),
		)
		os.Exit(1)
	}

	preflight, err := verify.CheckArchTrustPreflight(
		*pacmanKeyBinary,
		*archTrustDB,
		*archKeyringSource,
		nil,
	)
	if err != nil {
		ui.PrintVerificationPreflightBlocked(os.Stdout, err)
		os.Exit(1)
	}

	ui.PrintVerificationPreflight(os.Stdout, preflight)

	verifyCtx, cancelVerify := context.WithTimeout(
		context.Background(),
		60*time.Second,
	)
	defer cancelVerify()

	fetcher := fetch.New(nil, 8<<20)
	archVerifier := verify.NewPacmanKey()
	archVerifier.Binary = *pacmanKeyBinary

	results, err := verify.VerifyTarget(
		verifyCtx,
		target,
		*workDir,
		fetcher,
		archVerifier,
	)
	if err != nil {
		ui.PrintVerificationBlocked(os.Stdout, err)
		os.Exit(1)
	}

	ui.PrintVerificationResults(os.Stdout, results)
}
