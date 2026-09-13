package ui

import (
	"fmt"
	"io"

	"github.com/MikereDD/Reforge/internal/catalog"
)

func PrintInstallerMenu(w io.Writer, c *catalog.Catalog) {
	fmt.Fprintln(w, "REFORGE")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Install Operating System")
	fmt.Fprintln(w)

	for _, entry := range c.Entries {
		if entry.Kind != "os" {
			continue
		}

		fmt.Fprintf(w, "  %d. %s\n", entry.MenuOrder, entry.Name)
	}
}
