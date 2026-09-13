package verify

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

type fakeRunner struct {
	name string
	args []string
}

func (f *fakeRunner) Run(
	ctx context.Context,
	name string,
	args ...string,
) ([]byte, error) {
	f.name = name
	f.args = append([]string(nil), args...)

	return []byte("Good signature"), nil
}

func TestGPGVUsesExplicitTrustedKeyring(t *testing.T) {
	dir := t.TempDir()

	keyring := filepath.Join(dir, "archlinux.gpg")
	artifact := filepath.Join(dir, "ipxe-arch.efi")
	signature := filepath.Join(dir, "ipxe-arch.efi.sig")

	for _, path := range []string{
		keyring,
		artifact,
		signature,
	} {
		if err := os.WriteFile(
			path,
			[]byte("test"),
			0o600,
		); err != nil {
			t.Fatal(err)
		}
	}

	runner := &fakeRunner{}

	verifier := &GPGV{
		Binary:  "gpgv-test",
		Keyring: keyring,
		runner:  runner,
	}

	_, err := verifier.Verify(
		context.Background(),
		artifact,
		signature,
	)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}

	if runner.name != "gpgv-test" {
		t.Fatalf(
			"expected gpgv-test, got %q",
			runner.name,
		)
	}

	if len(runner.args) != 4 {
		t.Fatalf(
			"expected 4 gpgv arguments, got %d",
			len(runner.args),
		)
	}

	if runner.args[0] != "--keyring" ||
		runner.args[1] != keyring ||
		runner.args[2] != signature ||
		runner.args[3] != artifact {
		t.Fatalf(
			"unexpected gpgv arguments: %#v",
			runner.args,
		)
	}
}
