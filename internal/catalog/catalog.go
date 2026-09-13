package catalog

import "time"

type Catalog struct {
	SchemaVersion string    `json:"schemaVersion"`
	GeneratedAt   time.Time `json:"generatedAt"`
	Entries       []Entry   `json:"entries"`
}

type Entry struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Kind           string        `json:"kind"`
	MenuOrder      int           `json:"menuOrder"`
	Architecture   []string      `json:"architecture"`
	Channel        string        `json:"channel"`
	SourcePolicy   string        `json:"sourcePolicy"`
	BootMethod     string        `json:"bootMethod"`
	OnlineRequired bool          `json:"onlineRequired"`
	Recommended    bool          `json:"recommended"`
	Resolver       *Resolver     `json:"resolver,omitempty"`
	Verification   *Verification `json:"verification,omitempty"`
	Notes          []string      `json:"notes,omitempty"`
}

type Resolver struct {
	Type     string `json:"type"`
	Provider string `json:"provider,omitempty"`
	URL      string `json:"url,omitempty"`
}

type Verification struct {
	Required bool     `json:"required"`
	Methods  []string `json:"methods"`
}
