package ui

import (
	"fmt"
	"io"

	"github.com/MikereDD/Reforge/internal/verify"
)

func PrintVerificationPreflight(
	w io.Writer,
	result verify.PreflightResult,
) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Verification Preflight")
	fmt.Fprintln(w)

	fmt.Fprintf(
		w,
		"pacman-key:        PASS (%s)\n",
		result.VerifierPath,
	)

	fmt.Fprintf(
		w,
		"Pacman trust DB:   PASS (%s)\n",
		result.TrustDB,
	)

	fmt.Fprintf(
		w,
		"Arch keyring:      PASS (%s)\n",
		result.KeyringSource,
	)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Status: READY")
}

func PrintVerificationPreflightBlocked(
	w io.Writer,
	err error,
) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Verification Preflight")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Status: BLOCKED")
	fmt.Fprintf(w, "Reason: %v\n", err)
}
