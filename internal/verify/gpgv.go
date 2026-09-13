package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type commandRunner interface {
	Run(
		ctx context.Context,
		name string,
		args ...string,
	) ([]byte, error)
}

type execRunner struct{}

func (execRunner) Run(
	ctx context.Context,
	name string,
	args ...string,
) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	return cmd.CombinedOutput()
}

type GPGV struct {
	Binary  string
	Keyring string
	runner  commandRunner
}

func NewGPGV(keyring string) *GPGV {
	return &GPGV{
		Binary:  "gpgv",
		Keyring: keyring,
		runner:  execRunner{},
	}
}

func (v *GPGV) Verify(
	ctx context.Context,
	artifactPath string,
	signaturePath string,
) (string, error) {
	if v.Keyring == "" {
		return "", fmt.Errorf("trusted Arch keyring was not configured")
	}

	if _, err := os.Stat(v.Keyring); err != nil {
		return "", fmt.Errorf(
			"trusted Arch keyring unavailable: %w",
			err,
		)
	}

	if _, err := os.Stat(artifactPath); err != nil {
		return "", fmt.Errorf("artifact unavailable: %w", err)
	}

	if _, err := os.Stat(signaturePath); err != nil {
		return "", fmt.Errorf("signature unavailable: %w", err)
	}

	binary := v.Binary
	if binary == "" {
		binary = "gpgv"
	}

	runner := v.runner
	if runner == nil {
		runner = execRunner{}
	}

	output, err := runner.Run(
		ctx,
		binary,
		"--keyring",
		v.Keyring,
		signaturePath,
		artifactPath,
	)

	message := strings.TrimSpace(string(output))

	if err != nil {
		if message == "" {
			return "", fmt.Errorf("PGP verification failed: %w", err)
		}

		return message, fmt.Errorf(
			"PGP verification failed: %w: %s",
			err,
			message,
		)
	}

	return message, nil
}
