package verify

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type PacmanKey struct {
	Binary string
	runner commandRunner
}

func NewPacmanKey() *PacmanKey {
	return &PacmanKey{
		Binary: "pacman-key",
		runner: execRunner{},
	}
}

func (v *PacmanKey) Verify(
	ctx context.Context,
	artifactPath string,
	signaturePath string,
) (string, error) {
	if _, err := os.Stat(artifactPath); err != nil {
		return "", fmt.Errorf("artifact unavailable: %w", err)
	}

	if _, err := os.Stat(signaturePath); err != nil {
		return "", fmt.Errorf("signature unavailable: %w", err)
	}

	binary := v.Binary
	if binary == "" {
		binary = "pacman-key"
	}

	runner := v.runner
	if runner == nil {
		runner = execRunner{}
	}

	output, err := runner.Run(
		ctx,
		binary,
		"-v",
		signaturePath,
		artifactPath,
	)

	message := strings.TrimSpace(string(output))

	if err != nil {
		if message == "" {
			return "", fmt.Errorf(
				"Arch PGP verification failed: %w",
				err,
			)
		}

		return message, fmt.Errorf(
			"Arch PGP verification failed: %w: %s",
			err,
			message,
		)
	}

	return message, nil
}
