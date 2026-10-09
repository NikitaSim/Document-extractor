package service

import (
	"context"
	"fmt"
	"strings"

	"document-extractor/internal/ai"
	"document-extractor/internal/model"
	"document-extractor/internal/parser"
)

type DocumentService struct {
	parsers *parser.Registry
	ai      ai.Extractor
}

func NewDocumentService(parsers *parser.Registry, extractor ai.Extractor) *DocumentService {
	return &DocumentService{parsers: parsers, ai: extractor}
}

func (s *DocumentService) Extract(ctx context.Context, filename string, data []byte) (model.ExtractedData, error) {
	p, err := s.parsers.Get(filename)
	if err != nil {
		return model.ExtractedData{}, err
	}
	doc, err := p.Parse(data)
	if err != nil {
		return model.ExtractedData{}, err
	}
	if strings.TrimSpace(doc.Text) == "" {
		return model.ExtractedData{}, fmt.Errorf("document contains no extractable text")
	}
	return s.ai.Extract(ctx, doc)
}
