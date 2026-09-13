package fetch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var ErrInsecureURL = errors.New("only HTTPS artifact URLs are allowed")

type Result struct {
	URL    string
	Path   string
	Size   int64
	SHA256 string
}

type Fetcher struct {
	client   *http.Client
	maxBytes int64
}

func New(client *http.Client, maxBytes int64) *Fetcher {
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	if maxBytes <= 0 {
		maxBytes = 8 << 20
	}

	cloned := *client
	previousRedirect := cloned.CheckRedirect

	cloned.CheckRedirect = func(
		req *http.Request,
		via []*http.Request,
	) error {
		if req.URL.Scheme != "https" {
			return ErrInsecureURL
		}

		if previousRedirect != nil {
			return previousRedirect(req, via)
		}

		if len(via) >= 10 {
			return errors.New("too many HTTP redirects")
		}

		return nil
	}

	return &Fetcher{
		client:   &cloned,
		maxBytes: maxBytes,
	}
}

func (f *Fetcher) Fetch(
	ctx context.Context,
	rawURL string,
	destination string,
) (Result, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return Result{}, fmt.Errorf("parse artifact URL: %w", err)
	}

	if parsed.Scheme != "https" || parsed.Host == "" {
		return Result{}, ErrInsecureURL
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		rawURL,
		nil,
	)
	if err != nil {
		return Result{}, fmt.Errorf("create artifact request: %w", err)
	}

	req.Header.Set("User-Agent", "Reforge/0.2-dev")

	resp, err := f.client.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("download artifact: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf(
			"artifact server returned HTTP %d",
			resp.StatusCode,
		)
	}

	if resp.ContentLength > f.maxBytes {
		return Result{}, fmt.Errorf(
			"artifact exceeds maximum size of %d bytes",
			f.maxBytes,
		)
	}

	if err := os.MkdirAll(
		filepath.Dir(destination),
		0o755,
	); err != nil {
		return Result{}, fmt.Errorf("create artifact directory: %w", err)
	}

	partial := destination + ".part"

	file, err := os.OpenFile(
		partial,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return Result{}, fmt.Errorf("create artifact file: %w", err)
	}

	defer os.Remove(partial)

	hash := sha256.New()

	written, copyErr := io.Copy(
		io.MultiWriter(file, hash),
		io.LimitReader(resp.Body, f.maxBytes+1),
	)

	closeErr := file.Close()

	if copyErr != nil {
		return Result{}, fmt.Errorf("write artifact: %w", copyErr)
	}

	if closeErr != nil {
		return Result{}, fmt.Errorf("close artifact: %w", closeErr)
	}

	if written > f.maxBytes {
		return Result{}, fmt.Errorf(
			"artifact exceeds maximum size of %d bytes",
			f.maxBytes,
		)
	}

	if err := os.Remove(destination); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return Result{}, fmt.Errorf(
			"replace previous artifact: %w",
			err,
		)
	}

	if err := os.Rename(partial, destination); err != nil {
		return Result{}, fmt.Errorf("finalize artifact: %w", err)
	}

	return Result{
		URL:    rawURL,
		Path:   destination,
		Size:   written,
		SHA256: hex.EncodeToString(hash.Sum(nil)),
	}, nil
}
