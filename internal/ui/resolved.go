package ui

import (
	"fmt"
	"io"

	"github.com/MikereDD/Reforge/internal/resolver"
)

func PrintResolvedTarget(w io.Writer, target resolver.Target) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Resolved Installation Target")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "Operating system:  %s\n", target.Name)
	fmt.Fprintf(w, "Latest release:    %s\n", target.Version)
	fmt.Fprintf(w, "Kernel:            %s\n", target.KernelVersion)
	fmt.Fprintf(w, "Release date:      %s\n", target.ReleaseDate)
	fmt.Fprintf(w, "Provider:          %s\n", target.Provider)
	fmt.Fprintf(w, "Boot method:       %s\n", target.BootMethod)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Boot artifacts:")

	for _, artifact := range target.Artifacts {
		fmt.Fprintf(
			w,
			"  %s / %s\n",
			artifact.Firmware,
			artifact.Architecture,
		)

		fmt.Fprintf(w, "    %s\n", artifact.URL)
		fmt.Fprintf(w, "    Signature: %s\n", artifact.SignatureURL)
	}

	if len(target.Requirements) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Requirements:")

		for _, requirement := range target.Requirements {
			fmt.Fprintf(w, "  - %s\n", requirement)
		}
	}

	if len(target.Warnings) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Warnings:")

		for _, warning := range target.Warnings {
			fmt.Fprintf(w, "  - %s\n", warning)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Status: RESOLVED")
}
