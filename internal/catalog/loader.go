package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

func Load(path string) (*Catalog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open catalog: %w", err)
	}
	defer file.Close()

	return Decode(file)
}

func Decode(r io.Reader) (*Catalog, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var c Catalog

	if err := decoder.Decode(&c); err != nil {
		return nil, fmt.Errorf("decode catalog: %w", err)
	}

	if err := Validate(&c); err != nil {
		return nil, err
	}

	sort.SliceStable(c.Entries, func(i, j int) bool {
		return c.Entries[i].MenuOrder < c.Entries[j].MenuOrder
	})

	return &c, nil
}
