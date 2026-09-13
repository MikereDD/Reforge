package catalog

import (
	"fmt"
	"regexp"
)

var validID = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

var validKinds = map[string]bool{
	"os":         true,
	"rescue":     true,
	"diagnostic": true,
}

var validArchitectures = map[string]bool{
	"x86_64":  true,
	"aarch64": true,
}

var validSourcePolicies = map[string]bool{
	"official-only":      true,
	"official-preferred": true,
	"reforge-hosted":     true,
}

var validBootMethods = map[string]bool{
	"netboot": true,
	"wim":     true,
	"iso":     true,
	"uki":     true,
	"custom":  true,
}

var validVerificationMethods = map[string]bool{
	"sha256":       true,
	"sha512":       true,
	"gpg":          true,
	"authenticode": true,
}

func Validate(c *Catalog) error {
	if c == nil {
		return fmt.Errorf("catalog is nil")
	}

	if c.SchemaVersion == "" {
		return fmt.Errorf("catalog schemaVersion is required")
	}

	if c.GeneratedAt.IsZero() {
		return fmt.Errorf("catalog generatedAt is required")
	}

	if len(c.Entries) == 0 {
		return fmt.Errorf("catalog contains no entries")
	}

	ids := make(map[string]bool)
	menuOrders := make(map[int]bool)

	for i, entry := range c.Entries {
		if err := validateEntry(entry); err != nil {
			return fmt.Errorf("entry %d (%q): %w", i, entry.ID, err)
		}

		if ids[entry.ID] {
			return fmt.Errorf("duplicate entry id %q", entry.ID)
		}
		ids[entry.ID] = true

		if menuOrders[entry.MenuOrder] {
			return fmt.Errorf("duplicate menuOrder %d", entry.MenuOrder)
		}
		menuOrders[entry.MenuOrder] = true
	}

	return nil
}

func validateEntry(e Entry) error {
	if !validID.MatchString(e.ID) {
		return fmt.Errorf("invalid id %q", e.ID)
	}

	if e.Name == "" {
		return fmt.Errorf("name is required")
	}

	if !validKinds[e.Kind] {
		return fmt.Errorf("invalid kind %q", e.Kind)
	}

	if e.MenuOrder < 1 {
		return fmt.Errorf("menuOrder must be greater than zero")
	}

	if len(e.Architecture) == 0 {
		return fmt.Errorf("at least one architecture is required")
	}

	seenArchitectures := make(map[string]bool)

	for _, architecture := range e.Architecture {
		if !validArchitectures[architecture] {
			return fmt.Errorf("invalid architecture %q", architecture)
		}

		if seenArchitectures[architecture] {
			return fmt.Errorf("duplicate architecture %q", architecture)
		}

		seenArchitectures[architecture] = true
	}

	if e.Channel == "" {
		return fmt.Errorf("channel is required")
	}

	if !validSourcePolicies[e.SourcePolicy] {
		return fmt.Errorf("invalid sourcePolicy %q", e.SourcePolicy)
	}

	if !validBootMethods[e.BootMethod] {
		return fmt.Errorf("invalid bootMethod %q", e.BootMethod)
	}

	if e.Resolver != nil && e.Resolver.Type == "" {
		return fmt.Errorf("resolver type is required")
	}

	if e.Verification != nil {
		seen := make(map[string]bool)

		for _, method := range e.Verification.Methods {
			if !validVerificationMethods[method] {
				return fmt.Errorf("invalid verification method %q", method)
			}

			if seen[method] {
				return fmt.Errorf("duplicate verification method %q", method)
			}

			seen[method] = true
		}
	}

	return nil
}
