package resolver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MikereDD/Reforge/internal/catalog"
)

var ErrUnsupported = errors.New("resolver not implemented")

type Artifact struct {
	Name         string
	Firmware     string
	Architecture string
	URL          string
	SignatureURL string
}

type Target struct {
	ID            string
	Name          string
	Provider      string
	Version       string
	KernelVersion string
	ReleaseDate   string
	BootMethod    string
	Artifacts     []Artifact
	Requirements  []string
	Warnings      []string
}

type Registry struct {
	client *http.Client
}

func New(client *http.Client) *Registry {
	if client == nil {
		client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	return &Registry{
		client: client,
	}
}

func (r *Registry) Resolve(
	ctx context.Context,
	entry catalog.Entry,
) (Target, error) {
	if entry.Resolver == nil {
		return Target{}, fmt.Errorf(
			"%w: %s has no resolver",
			ErrUnsupported,
			entry.Name,
		)
	}

	switch entry.Resolver.Provider {
	case "archlinux":
		return r.resolveArch(ctx, entry)

	default:
		return Target{}, fmt.Errorf(
			"%w: provider %q",
			ErrUnsupported,
			entry.Resolver.Provider,
		)
	}
}
