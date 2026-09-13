package fetch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestFetchHTTPSArtifact(t *testing.T) {
	const payload = "reforge-artifact"

	server := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(payload))
		}),
	)
	defer server.Close()

	fetcher := New(server.Client(), 1024)

	destination := filepath.Join(
		t.TempDir(),
		"artifact.bin",
	)

	result, err := fetcher.Fetch(
		context.Background(),
		server.URL,
		destination,
	)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	expected := sha256.Sum256([]byte(payload))
	expectedHash := hex.EncodeToString(expected[:])

	if result.SHA256 != expectedHash {
		t.Fatalf(
			"expected SHA-256 %q, got %q",
			expectedHash,
			result.SHA256,
		)
	}
}

func TestFetchRejectsHTTP(t *testing.T) {
	fetcher := New(nil, 1024)

	_, err := fetcher.Fetch(
		context.Background(),
		"http://example.invalid/artifact",
		filepath.Join(t.TempDir(), "artifact"),
	)

	if !errors.Is(err, ErrInsecureURL) {
		t.Fatalf(
			"expected ErrInsecureURL, got %v",
			err,
		)
	}
}
