package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MikereDD/Reforge/internal/catalog"
)

const (
	archNetbootUEFIURL = "https://archlinux.org/static/netboot/ipxe-arch.efi"

	archNetbootUEFISignatureURL = "https://archlinux.org/static/netboot/ipxe-arch.efi.sig"

	archNetbootBIOSURL = "https://archlinux.org/static/netboot/ipxe-arch.lkrn"

	archNetbootBIOSSignatureURL = "https://archlinux.org/static/netboot/ipxe-arch.lkrn.sig"
)

type archReleaseFeed struct {
	Releases      []archRelease `json:"releases"`
	LatestVersion string        `json:"latest_version"`
}

type archRelease struct {
	ReleaseDate   string `json:"release_date"`
	Version       string `json:"version"`
	KernelVersion string `json:"kernel_version"`
	Available     bool   `json:"available"`
}

func (r *Registry) resolveArch(
	ctx context.Context,
	entry catalog.Entry,
) (Target, error) {
	if entry.Resolver == nil {
		return Target{}, fmt.Errorf("Arch resolver metadata is missing")
	}

	if entry.Resolver.Type != "upstream-latest" {
		return Target{}, fmt.Errorf(
			"unsupported Arch resolver type %q",
			entry.Resolver.Type,
		)
	}

	if entry.Resolver.URL == "" {
		return Target{}, fmt.Errorf("Arch release API URL is missing")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		entry.Resolver.URL,
		nil,
	)
	if err != nil {
		return Target{}, fmt.Errorf("create Arch release request: %w", err)
	}

	req.Header.Set("User-Agent", "Reforge/0.2-dev")

	resp, err := r.client.Do(req)
	if err != nil {
		return Target{}, fmt.Errorf("query Arch release API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Target{}, fmt.Errorf(
			"Arch release API returned HTTP %d",
			resp.StatusCode,
		)
	}

	var feed archReleaseFeed

	decoder := json.NewDecoder(
		io.LimitReader(resp.Body, 16<<20),
	)

	if err := decoder.Decode(&feed); err != nil {
		return Target{}, fmt.Errorf(
			"decode Arch release API response: %w",
			err,
		)
	}

	if feed.LatestVersion == "" {
		return Target{}, fmt.Errorf(
			"Arch release API did not provide latest_version",
		)
	}

	var latest *archRelease

	for i := range feed.Releases {
		release := &feed.Releases[i]

		if release.Version == feed.LatestVersion {
			latest = release
			break
		}
	}

	if latest == nil {
		return Target{}, fmt.Errorf(
			"Arch latest release %q was not present in release list",
			feed.LatestVersion,
		)
	}

	if !latest.Available {
		return Target{}, fmt.Errorf(
			"Arch latest release %q is not currently available",
			latest.Version,
		)
	}

	return Target{
		ID:            entry.ID,
		Name:          entry.Name,
		Provider:      "archlinux",
		Version:       latest.Version,
		KernelVersion: latest.KernelVersion,
		ReleaseDate:   latest.ReleaseDate,
		BootMethod:    entry.BootMethod,

		Artifacts: []Artifact{
			{
				Name:         "Arch Linux Netboot",
				Firmware:     "UEFI",
				Architecture: "x86_64",
				URL:          archNetbootUEFIURL,
				SignatureURL: archNetbootUEFISignatureURL,
			},
			{
				Name:         "Arch Linux Netboot",
				Firmware:     "BIOS",
				Architecture: "x86_64",
				URL:          archNetbootBIOSURL,
				SignatureURL: archNetbootBIOSSignatureURL,
			},
		},

		Requirements: []string{
			"Wired Ethernet with DHCP auto-configuration",
			"Sufficient RAM to download and run the Arch live system",
		},

		Warnings: []string{
			"Secure Boot must be disabled for the current Arch UEFI netboot path",
		},
	}, nil
}
