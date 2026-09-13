package verify

import (
	"fmt"
	"os"
	"os/exec"
)

type ToolLookup func(string) (string, error)

type PreflightResult struct {
	Verifier      string
	VerifierPath  string
	TrustDB       string
	KeyringSource string
}

func CheckArchTrustPreflight(
	binary string,
	trustDB string,
	keyringSource string,
	lookup ToolLookup,
) (PreflightResult, error) {
	if binary == "" {
		binary = "pacman-key"
	}

	if trustDB == "" {
		return PreflightResult{}, fmt.Errorf(
			"Arch pacman trust database was not configured",
		)
	}

	trustInfo, err := os.Stat(trustDB)
	if err != nil {
		return PreflightResult{}, fmt.Errorf(
			"Arch pacman trust database unavailable: %w",
			err,
		)
	}

	if !trustInfo.IsDir() {
		return PreflightResult{}, fmt.Errorf(
			"Arch pacman trust database is not a directory: %s",
			trustDB,
		)
	}

	if keyringSource == "" {
		return PreflightResult{}, fmt.Errorf(
			"Arch keyring source was not configured",
		)
	}

	keyringInfo, err := os.Stat(keyringSource)
	if err != nil {
		return PreflightResult{}, fmt.Errorf(
			"Arch keyring source unavailable: %w",
			err,
		)
	}

	if keyringInfo.IsDir() {
		return PreflightResult{}, fmt.Errorf(
			"Arch keyring source is a directory: %s",
			keyringSource,
		)
	}

	if lookup == nil {
		lookup = exec.LookPath
	}

	binaryPath, err := lookup(binary)
	if err != nil {
		return PreflightResult{}, fmt.Errorf(
			"%s is unavailable: %w",
			binary,
			err,
		)
	}

	return PreflightResult{
		Verifier:      binary,
		VerifierPath:  binaryPath,
		TrustDB:       trustDB,
		KeyringSource: keyringSource,
	}, nil
}
