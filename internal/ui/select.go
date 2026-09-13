package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/MikereDD/Reforge/internal/catalog"
)

var ErrCanceled = errors.New("selection canceled")

func SelectInstaller(r *bufio.Reader, w io.Writer, c *catalog.Catalog) (catalog.Entry, error) {
	for {
		fmt.Fprint(w, "\nSelect an operating system [q to quit]: ")

		input, err := r.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return catalog.Entry{}, fmt.Errorf("read selection: %w", err)
		}

		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "q") {
			return catalog.Entry{}, ErrCanceled
		}

		selection, parseErr := strconv.Atoi(input)
		if parseErr != nil {
			fmt.Fprintln(w, "Invalid selection.")
			if errors.Is(err, io.EOF) {
				return catalog.Entry{}, io.EOF
			}
			continue
		}

		for _, entry := range c.Entries {
			if entry.Kind == "os" && entry.MenuOrder == selection {
				return entry, nil
			}
		}

		fmt.Fprintln(w, "Invalid selection.")

		if errors.Is(err, io.EOF) {
			return catalog.Entry{}, io.EOF
		}
	}
}

func PrintInstallerDetails(w io.Writer, entry catalog.Entry) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, entry.Name)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "Channel:          %s\n", entry.Channel)
	fmt.Fprintf(w, "Architecture:     %s\n", strings.Join(entry.Architecture, ", "))
	fmt.Fprintf(w, "Source policy:    %s\n", entry.SourcePolicy)
	fmt.Fprintf(w, "Boot method:      %s\n", entry.BootMethod)

	if entry.OnlineRequired {
		fmt.Fprintln(w, "Online required:  yes")
	} else {
		fmt.Fprintln(w, "Online required:  no")
	}

	if entry.Verification != nil && len(entry.Verification.Methods) > 0 {
		fmt.Fprintf(
			w,
			"Verification:     %s\n",
			strings.Join(entry.Verification.Methods, ", "),
		)
	} else {
		fmt.Fprintln(w, "Verification:     none")
	}
}

func Confirm(r *bufio.Reader, w io.Writer) (bool, error) {
	fmt.Fprint(w, "\nProceed? [y/N]: ")

	input, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("read confirmation: %w", err)
	}

	switch strings.ToLower(strings.TrimSpace(input)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
