package ui

import (
	"fmt"
	"io"

	"github.com/MikereDD/Reforge/internal/verify"
)

func PrintVerificationResults(
	w io.Writer,
	results []verify.ArtifactResult,
) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Artifact Verification")
	fmt.Fprintln(w)

	for _, result := range results {
		fmt.Fprintf(
			w,
			"%s / %s\n",
			result.Firmware,
			result.Architecture,
		)

		fmt.Fprintf(
			w,
			"  Download:       PASS (%d bytes)\n",
			result.Size,
		)

		fmt.Fprintf(
			w,
			"  SHA-256:        %s\n",
			result.SHA256,
		)

		fmt.Fprintln(
			w,
			"  PGP signature:  VALID",
		)

		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "Status: VERIFIED")
}

func PrintVerificationBlocked(
	w io.Writer,
	err error,
) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Artifact Verification")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Status: BLOCKED")
	fmt.Fprintf(w, "Reason: %v\n", err)
}
