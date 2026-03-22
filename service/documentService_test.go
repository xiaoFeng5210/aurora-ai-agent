package service

import (
	"aurora-agent/handler/dto"
	"testing"
)

func TestNormalizeDocumentDisplayName(t *testing.T) {
	if _, err := normalizeDocumentDisplayName(""); err != ErrDocumentDisplayNameRequired {
		t.Fatalf("expected %v, got %v", ErrDocumentDisplayNameRequired, err)
	}

	if _, err := normalizeDocumentDisplayName("   "); err != ErrDocumentDisplayNameRequired {
		t.Fatalf("expected %v, got %v", ErrDocumentDisplayNameRequired, err)
	}

	displayName, err := normalizeDocumentDisplayName("  Aurora Doc  ")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if displayName != "Aurora Doc" {
		t.Fatalf("unexpected display name: %q", displayName)
	}
}

func TestNormalizeOptionalString(t *testing.T) {
	if normalizeOptionalString(nil) != nil {
		t.Fatal("expected nil when input is nil")
	}

	blank := "   "
	if normalizeOptionalString(&blank) != nil {
		t.Fatal("expected nil for blank string")
	}

	value := "  file.txt  "
	normalized := normalizeOptionalString(&value)
	if normalized == nil {
		t.Fatal("expected normalized value, got nil")
	}
	if *normalized != "file.txt" {
		t.Fatalf("unexpected normalized value: %q", *normalized)
	}
}

func TestBuildDocumentUpdates(t *testing.T) {
	if _, err := buildDocumentUpdates(dto.UpdateDocumentRequest{}); err != ErrNoFieldsToUpdate {
		t.Fatalf("expected %v, got %v", ErrNoFieldsToUpdate, err)
	}

	blankDisplayName := "   "
	if _, err := buildDocumentUpdates(dto.UpdateDocumentRequest{DisplayName: &blankDisplayName}); err != ErrDocumentDisplayNameRequired {
		t.Fatalf("expected %v, got %v", ErrDocumentDisplayNameRequired, err)
	}

	fileName := "   "
	updates, err := buildDocumentUpdates(dto.UpdateDocumentRequest{FileName: &fileName})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value, ok := updates["file_name"]; !ok || value != nil {
		t.Fatalf("expected file_name to be normalized to nil, got %#v", updates["file_name"])
	}
}

func TestNormalizePagination(t *testing.T) {
	page, pageSize := normalizePagination(0, 0)
	if page != 1 {
		t.Fatalf("expected page 1, got %d", page)
	}
	if pageSize != defaultPageSize {
		t.Fatalf("expected default page size %d, got %d", defaultPageSize, pageSize)
	}

	page, pageSize = normalizePagination(2, maxPageSize+1)
	if page != 2 {
		t.Fatalf("expected page 2, got %d", page)
	}
	if pageSize != maxPageSize {
		t.Fatalf("expected max page size %d, got %d", maxPageSize, pageSize)
	}
}
