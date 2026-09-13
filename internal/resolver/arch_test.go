package resolver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MikereDD/Reforge/internal/catalog"
)

func TestResolveArch(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			fmt.Fprint(w, `{
"releases": [
{
"release_date": "2026-09-01",
"version": "2026.09.01",
"kernel_version": "7.2.2",
"available": true
}
],
"latest_version": "2026.09.01"
}`)
		}),
	)
	defer server.Close()

	entry := catalog.Entry{
		ID:         "arch-linux",
		Name:       "Arch Linux",
		BootMethod: "netboot",
		Resolver: &catalog.Resolver{
			Type:     "upstream-latest",
			Provider: "archlinux",
			URL:      server.URL,
		},
	}

	registry := New(server.Client())

	target, err := registry.Resolve(context.Background(), entry)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if target.Version != "2026.09.01" {
		t.Fatalf(
			"expected version 2026.09.01, got %q",
			target.Version,
		)
	}

	if target.KernelVersion != "7.2.2" {
		t.Fatalf(
			"expected kernel 7.2.2, got %q",
			target.KernelVersion,
		)
	}

	if len(target.Artifacts) != 2 {
		t.Fatalf(
			"expected 2 netboot artifacts, got %d",
			len(target.Artifacts),
		)
	}

	if target.Artifacts[0].URL != archNetbootUEFIURL {
		t.Fatalf(
			"unexpected UEFI artifact URL %q",
			target.Artifacts[0].URL,
		)
	}
}

func TestResolveArchRejectsUnavailableLatest(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			fmt.Fprint(w, `{
"releases": [
{
"release_date": "2026-09-01",
"version": "2026.09.01",
"kernel_version": "7.2.2",
"available": false
}
],
"latest_version": "2026.09.01"
}`)
		}),
	)
	defer server.Close()

	entry := catalog.Entry{
		ID:         "arch-linux",
		Name:       "Arch Linux",
		BootMethod: "netboot",
		Resolver: &catalog.Resolver{
			Type:     "upstream-latest",
			Provider: "archlinux",
			URL:      server.URL,
		},
	}

	registry := New(server.Client())

	_, err := registry.Resolve(context.Background(), entry)

	if err == nil {
		t.Fatal("expected unavailable release to fail")
	}
}
