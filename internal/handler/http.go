package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"document-extractor/internal/service"
)

type Handler struct {
	service     *service.DocumentService
	maxFileSize int64
}

func New(s *service.DocumentService, maxFileSize int64) *Handler {
	return &Handler{service: s, maxFileSize: maxFileSize}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *Handler) Extract(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxFileSize+1024*1024)
	if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_MULTIPART", "invalid multipart request")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "FILE_REQUIRED", "field 'file' is required")
		return
	}
	defer file.Close()

	if header.Size > h.maxFileSize {
		writeError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", fmt.Sprintf("maximum file size is %d bytes", h.maxFileSize))
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" && ext != ".docx" {
		writeError(w, http.StatusUnsupportedMediaType, "UNSUPPORTED_FORMAT", "only PDF and DOCX are supported")
		return
	}

	data := make([]byte, header.Size)
	if _, err := io.ReadFull(file, data); err != nil {
		writeError(w, http.StatusBadRequest, "FILE_READ_ERROR", "failed to read uploaded file")
		return
	}

	ctx, cancel := contextWithTimeout(r, 5*time.Minute)
	defer cancel()
	result, err := h.service.Extract(ctx, header.Filename, data)
	if err != nil {
		status := http.StatusInternalServerError
		code := "EXTRACTION_FAILED"
		if strings.Contains(err.Error(), "unsupported") || strings.Contains(err.Error(), "no extractable") {
			status = http.StatusUnprocessableEntity
			code = "DOCUMENT_INVALID"
		}
		writeError(w, status, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "filename": header.Filename, "data": result})
}

func contextWithTimeout(r *http.Request, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), d)
}

type errorResponse struct {
	Success bool      `json:"success"`
	Error   errorBody `json:"error"`
}
type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{false, errorBody{code, message}})
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
