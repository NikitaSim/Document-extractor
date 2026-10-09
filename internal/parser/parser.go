package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"document-extractor/internal/model"
	"github.com/ledongthuc/pdf"
)

type Parser interface {
	Parse(data []byte) (model.Document, error)
}

type Registry struct{ parsers map[string]Parser }

func NewRegistry() *Registry {
	return &Registry{parsers: map[string]Parser{
		".pdf":  PDFParser{},
		".docx": DOCXParser{},
	}}
}

func (r *Registry) Get(filename string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	p, ok := r.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported document format: %s", ext)
	}
	return p, nil
}

type PDFParser struct{}

func (PDFParser) Parse(data []byte) (model.Document, error) {
	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return model.Document{}, fmt.Errorf("open pdf: %w", err)
	}
	textReader, err := reader.GetPlainText()
	if err != nil {
		return model.Document{}, fmt.Errorf("extract pdf text: %w", err)
	}
	var textBuffer bytes.Buffer
	if _, err := textBuffer.ReadFrom(textReader); err != nil {
		return model.Document{}, fmt.Errorf("read extracted pdf text: %w", err)
	}
	return model.Document{Text: strings.TrimSpace(textBuffer.String()), Format: "pdf"}, nil
}

type DOCXParser struct{}

func (DOCXParser) Parse(data []byte) (model.Document, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return model.Document{}, fmt.Errorf("open docx: %w", err)
	}
	var documentXML io.ReadCloser
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			documentXML, err = f.Open()
			if err != nil {
				return model.Document{}, fmt.Errorf("open docx xml: %w", err)
			}
			break
		}
	}
	if documentXML == nil {
		return model.Document{}, fmt.Errorf("word/document.xml not found")
	}
	defer documentXML.Close()

	var b strings.Builder
	dec := xml.NewDecoder(documentXML)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return model.Document{}, fmt.Errorf("read docx xml: %w", err)
		}
		switch t := tok.(type) {
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				b.WriteString(text)
			}
		case xml.StartElement:
			if t.Name.Local == "tab" {
				b.WriteByte('\t')
			}
		case xml.EndElement:
			if t.Name.Local == "p" || t.Name.Local == "tr" {
				b.WriteByte('\n')
			}
		}
	}
	return model.Document{Text: strings.TrimSpace(b.String()), Format: "docx"}, nil
}
