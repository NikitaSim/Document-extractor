package parser

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	if _, err := r.Get("document.pdf"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Get("document.DOCX"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Get("document.txt"); err == nil {
		t.Fatal("expected unsupported format error")
	}
}

func TestDOCXParser(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write([]byte(`<?xml version="1.0"?><w:document xmlns:w="x"><w:body><w:p><w:r><w:t>Hello</w:t></w:r><w:r><w:t> world</w:t></w:r></w:p><w:p><w:r><w:t>Second</w:t></w:r></w:p></w:body></w:document>`))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	doc, err := (DOCXParser{}).Parse(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if doc.Format != "docx" {
		t.Fatalf("format = %q", doc.Format)
	}
	if doc.Text != "Helloworld\nSecond" {
		t.Fatalf("text = %q", doc.Text)
	}
}
