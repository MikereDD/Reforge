package verify

import (
	"fmt"
	"os"
	"os/exec"
)

type ToolLookup func(string) (string, error)

type PreflightResult struct {
	Binary     string
	BinaryPath string
	Keyring    string
}

func CheckGPGVPreflight(binary string, keyring string, lookup ToolLookup) (PreflightResult, error) {
	if binary == "" {
		binary = "gpgv"
	}

	if keyring == "" {
		return PreflightResult{}, fmt.Errorf("trusted Arch keyring was not configured")
	}

	info, err := os.Stat(keyring)
	if err != nil {
		return PreflightResult{}, fmt.Errorf("trusted Arch keyring unavailable: %w", err)
	}

	if info.IsDir() {
		return PreflightResult{}, fmt.Errorf("trusted Arch keyring path is a directory: %s", keyring)
	}

	if lookup == nil {
		lookup = exec.LookPath
	}

	binaryPath, err := lookup(binary)
	if err != nil {
		return PreflightResult{}, fmt.Errorf("%s is unavailable: %w", binary, err)
	}

	return PreflightResult{
		Binary:     binary,
		BinaryPath: binaryPath,
		Keyring:    keyring,
	}, nil
}
