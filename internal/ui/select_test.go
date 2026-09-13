package ui

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/MikereDD/Reforge/internal/catalog"
)

func testCatalog() *catalog.Catalog {
	return &catalog.Catalog{
		Entries: []catalog.Entry{
			{
				ID:           "windows-11",
				Name:         "Windows 11",
				Kind:         "os",
				MenuOrder:    1,
				Architecture: []string{"x86_64"},
			},
			{
				ID:           "arch-linux",
				Name:         "Arch Linux",
				Kind:         "os",
				MenuOrder:    2,
				Architecture: []string{"x86_64"},
			},
		},
	}
}

func TestSelectInstaller(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("2\n"))
	var output bytes.Buffer

	entry, err := SelectInstaller(reader, &output, testCatalog())
	if err != nil {
		t.Fatalf("SelectInstaller returned error: %v", err)
	}

	if entry.ID != "arch-linux" {
		t.Fatalf("expected arch-linux, got %q", entry.ID)
	}
}

func TestSelectInstallerRetriesInvalidSelection(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("99\n1\n"))
	var output bytes.Buffer

	entry, err := SelectInstaller(reader, &output, testCatalog())
	if err != nil {
		t.Fatalf("SelectInstaller returned error: %v", err)
	}

	if entry.ID != "windows-11" {
		t.Fatalf("expected windows-11, got %q", entry.ID)
	}

	if !strings.Contains(output.String(), "Invalid selection.") {
		t.Fatal("expected invalid selection message")
	}
}

func TestSelectInstallerCancel(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("q\n"))
	var output bytes.Buffer

	_, err := SelectInstaller(reader, &output, testCatalog())

	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("expected ErrCanceled, got %v", err)
	}
}

func TestConfirm(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("yes\n"))
	var output bytes.Buffer

	ok, err := Confirm(reader, &output)
	if err != nil {
		t.Fatalf("Confirm returned error: %v", err)
	}

	if !ok {
		t.Fatal("expected confirmation")
	}
}
