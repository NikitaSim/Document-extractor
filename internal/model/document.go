package model

type Document struct {
	Text   string
	Format string
}

type ExtractedData struct {
	DocumentType string         `json:"document_type"`
	Summary      string         `json:"summary"`
	Fields       map[string]any `json:"fields"`
}
