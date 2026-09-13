package verify

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

type pacmanKeyFakeRunner struct {
	name string
	args []string
}

func (f *pacmanKeyFakeRunner) Run(
	ctx context.Context,
	name string,
	args ...string,
) ([]byte, error) {
	f.name = name
	f.args = append([]string(nil), args...)

	return []byte(`gpg: Good signature from "Arch Signer <signer@archlinux.org>" [full]`), nil
}

func TestPacmanKeyVerifyUsesDetachedSignature(t *testing.T) {
	dir := t.TempDir()

	artifact := filepath.Join(dir, "ipxe-arch.efi")
	signature := filepath.Join(dir, "ipxe-arch.efi.sig")

	for _, path := range []string{
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

	runner := &pacmanKeyFakeRunner{}

	verifier := &PacmanKey{
		Binary: "pacman-key-test",
		runner: runner,
	}

	output, err := verifier.Verify(
		context.Background(),
		artifact,
		signature,
	)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}

	if runner.name != "pacman-key-test" {
		t.Fatalf(
			"expected pacman-key-test, got %q",
			runner.name,
		)
	}

	if len(runner.args) != 3 {
		t.Fatalf(
			"expected 3 arguments, got %d",
			len(runner.args),
		)
	}

	if runner.args[0] != "-v" ||
		runner.args[1] != signature ||
		runner.args[2] != artifact {
		t.Fatalf(
			"unexpected pacman-key arguments: %#v",
			runner.args,
		)
	}

	if output == "" {
		t.Fatal("expected verifier output")
	}
}
