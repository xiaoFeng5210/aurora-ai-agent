package embedding

import (
	"io"
	"os"
	"testing"
)

func TestMdDocument_ConvertToText(t *testing.T) {
	file, err := os.Open("test.md")
	if err != nil {
		t.Fatalf("Failed to open test.md: %v", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read test.md: %v", err)
	}

	md := &MdDocument{
		data: data,
	}

	md.ConvertToText()

	t.Logf("content: %s", md.content)
}

func TestMdDocumentEmbedding(t *testing.T) {
	file, err := os.Open("test.md")
	if err != nil {
		t.Fatalf("Failed to open test.md: %v", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read test.md: %v", err)
	}

	md := &MdDocument{
		data: data,
	}

	md.ConvertToText()
	chunks := md.Chunk()
	t.Logf("chunks: %v", chunks)
	vectors, err := md.Embedding()
	if err != nil {
		t.Fatalf("Failed to embed: %v", err)
	}
	t.Logf("vectors: %v", len(vectors))
	// updateResult, err := md.UpsertQdrantVector()
	// if err != nil {
	// 	t.Fatalf("Failed to upsert qdrant vector: %v", err)
	// }
	// t.Logf("updateResult: %v", updateResult)
}
