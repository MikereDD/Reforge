package verify

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func makeArchTrustFixture(t *testing.T) (string, string) {
	t.Helper()

	dir := t.TempDir()
	trustDB := filepath.Join(dir, "gnupg")
	keyring := filepath.Join(dir, "archlinux.gpg")

	if err := os.MkdirAll(trustDB, 0o700); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		keyring,
		[]byte("arch-keyring-source"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	return trustDB, keyring
}

func TestCheckArchTrustPreflight(t *testing.T) {
	trustDB, keyring := makeArchTrustFixture(t)

	result, err := CheckArchTrustPreflight(
		"pacman-key-test",
		trustDB,
		keyring,
		func(name string) (string, error) {
			if name != "pacman-key-test" {
				t.Fatalf("unexpected binary lookup %q", name)
			}

			return "/usr/bin/pacman-key-test", nil
		},
	)
	if err != nil {
		t.Fatalf(
			"CheckArchTrustPreflight returned error: %v",
			err,
		)
	}

	if result.VerifierPath != "/usr/bin/pacman-key-test" {
		t.Fatalf(
			"expected verifier path /usr/bin/pacman-key-test, got %q",
			result.VerifierPath,
		)
	}

	if result.TrustDB != trustDB {
		t.Fatalf(
			"expected trust DB %q, got %q",
			trustDB,
			result.TrustDB,
		)
	}

	if result.KeyringSource != keyring {
		t.Fatalf(
			"expected keyring source %q, got %q",
			keyring,
			result.KeyringSource,
		)
	}
}

func TestCheckArchTrustPreflightRejectsMissingTrustDB(t *testing.T) {
	_, keyring := makeArchTrustFixture(t)

	_, err := CheckArchTrustPreflight(
		"pacman-key",
		filepath.Join(t.TempDir(), "missing"),
		keyring,
		func(name string) (string, error) {
			return "/usr/bin/pacman-key", nil
		},
	)

	if err == nil {
		t.Fatal("expected missing trust DB to fail")
	}
}

func TestCheckArchTrustPreflightRejectsMissingBinary(t *testing.T) {
	trustDB, keyring := makeArchTrustFixture(t)

	_, err := CheckArchTrustPreflight(
		"pacman-key",
		trustDB,
		keyring,
		func(name string) (string, error) {
			return "", errors.New("not found")
		},
	)

	if err == nil {
		t.Fatal("expected missing binary to fail")
	}
}
