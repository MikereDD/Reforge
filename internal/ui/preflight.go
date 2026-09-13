package ui

import (
	"fmt"
	"io"

	"github.com/MikereDD/Reforge/internal/verify"
)

func PrintVerificationPreflight(w io.Writer, result verify.PreflightResult) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Verification Preflight")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "gpgv:             PASS (%s)\n", result.BinaryPath)
	fmt.Fprintf(w, "Trusted keyring:  PASS (%s)\n", result.Keyring)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Status: READY")
}

func PrintVerificationPreflightBlocked(w io.Writer, err error) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Verification Preflight")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Status: BLOCKED")
	fmt.Fprintf(w, "Reason: %v\n", err)
}
