package verify

import (
	"context"
	"fmt"
	"net/url"
	"os"
	pathpkg "path"
	"path/filepath"

	"github.com/MikereDD/Reforge/internal/fetch"
	"github.com/MikereDD/Reforge/internal/resolver"
)

type artifactFetcher interface {
	Fetch(
		ctx context.Context,
		rawURL string,
		destination string,
	) (fetch.Result, error)
}

type detachedVerifier interface {
	Verify(
		ctx context.Context,
		artifactPath string,
		signaturePath string,
	) (string, error)
}

type ArtifactResult struct {
	Firmware     string
	Architecture string
	Path         string
	Size         int64
	SHA256       string
	Signature    string
}

func VerifyTarget(
	ctx context.Context,
	target resolver.Target,
	workDir string,
	fetcher artifactFetcher,
	verifier detachedVerifier,
) ([]ArtifactResult, error) {
	if fetcher == nil {
		return nil, fmt.Errorf("artifact fetcher is required")
	}

	if verifier == nil {
		return nil, fmt.Errorf("signature verifier is required")
	}

	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return nil, fmt.Errorf(
			"create verification workspace: %w",
			err,
		)
	}

	results := make([]ArtifactResult, 0, len(target.Artifacts))

	for _, artifact := range target.Artifacts {
		if artifact.SignatureURL == "" {
			return nil, fmt.Errorf(
				"%s/%s has no signature URL",
				artifact.Firmware,
				artifact.Architecture,
			)
		}

		artifactName, err := URLFileName(artifact.URL)
		if err != nil {
			return nil, err
		}

		signatureName, err := URLFileName(artifact.SignatureURL)
		if err != nil {
			return nil, err
		}

		artifactPath := filepath.Join(workDir, artifactName)
		signaturePath := filepath.Join(workDir, signatureName)

		download, err := fetcher.Fetch(
			ctx,
			artifact.URL,
			artifactPath,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"fetch %s/%s artifact: %w",
				artifact.Firmware,
				artifact.Architecture,
				err,
			)
		}

		if _, err := fetcher.Fetch(
			ctx,
			artifact.SignatureURL,
			signaturePath,
		); err != nil {
			return nil, fmt.Errorf(
				"fetch %s/%s signature: %w",
				artifact.Firmware,
				artifact.Architecture,
				err,
			)
		}

		signatureOutput, err := verifier.Verify(
			ctx,
			artifactPath,
			signaturePath,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"verify %s/%s: %w",
				artifact.Firmware,
				artifact.Architecture,
				err,
			)
		}

		results = append(results, ArtifactResult{
			Firmware:     artifact.Firmware,
			Architecture: artifact.Architecture,
			Path:         download.Path,
			Size:         download.Size,
			SHA256:       download.SHA256,
			Signature:    signatureOutput,
		})
	}

	return results, nil
}

func URLFileName(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse artifact URL: %w", err)
	}

	name := pathpkg.Base(parsed.Path)

	if name == "" || name == "." || name == "/" {
		return "", fmt.Errorf(
			"artifact URL %q has no file name",
			rawURL,
		)
	}

	return name, nil
}
