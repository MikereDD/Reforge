package verify

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckGPGVPreflight(t *testing.T) {
	keyring := filepath.Join(t.TempDir(), "archlinux.gpg")
	if err := os.WriteFile(keyring, []byte("trusted-keyring"), 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := CheckGPGVPreflight(
		"gpgv-test",
		keyring,
		func(name string) (string, error) {
			if name != "gpgv-test" {
				t.Fatalf("unexpected binary lookup %q", name)
			}
			return "/usr/bin/gpgv-test", nil
		},
	)
	if err != nil {
		t.Fatalf("CheckGPGVPreflight returned error: %v", err)
	}

	if result.BinaryPath != "/usr/bin/gpgv-test" {
		t.Fatalf("expected binary path /usr/bin/gpgv-test, got %q", result.BinaryPath)
	}
	if result.Keyring != keyring {
		t.Fatalf("expected keyring %q, got %q", keyring, result.Keyring)
	}
}

func TestCheckGPGVPreflightRejectsMissingKeyring(t *testing.T) {
	_, err := CheckGPGVPreflight(
		"gpgv",
		"",
		func(name string) (string, error) {
			return "/usr/bin/gpgv", nil
		},
	)
	if err == nil {
		t.Fatal("expected missing keyring to fail")
	}
}

func TestCheckGPGVPreflightRejectsMissingBinary(t *testing.T) {
	keyring := filepath.Join(t.TempDir(), "archlinux.gpg")
	if err := os.WriteFile(keyring, []byte("trusted-keyring"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := CheckGPGVPreflight(
		"gpgv",
		keyring,
		func(name string) (string, error) {
			return "", errors.New("not found")
		},
	)
	if err == nil {
		t.Fatal("expected missing binary to fail")
	}
}
