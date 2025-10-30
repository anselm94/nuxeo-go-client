package nuxeo

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument("File", "myfile")
	if doc.Type != "File" {
		t.Errorf("expected Type 'File', got '%s'", doc.Type)
	}
	if doc.Name != "myfile" {
		t.Errorf("expected Name 'myfile', got '%s'", doc.Name)
	}
	if doc.entity.EntityType != EntityTypeDocument {
		t.Errorf("expected EntityType 'document', got '%s'", doc.entity.EntityType)
	}
	if _, ok := doc.Properties[DocumentPropertyDCTitle]; !ok {
		t.Errorf("expected dc:title property map, got none")
	}
}

func TestFacetMethods(t *testing.T) {
	doc := NewDocument("File", "myfile")
	doc.Facets = []string{"Folderish", "Collection"}

	t.Run("HasFacet", func(t *testing.T) {
		if !doc.HasFacet("Folderish") {
			t.Error("expected HasFacet to return true for 'Folderish'")
		}
		if doc.HasFacet("Nonexistent") {
			t.Error("expected HasFacet to return false for 'Nonexistent'")
		}
	})
	t.Run("IsFolder", func(t *testing.T) {
		if !doc.IsFolder() {
			t.Error("expected IsFolder to return true")
		}
	})
	t.Run("IsCollection", func(t *testing.T) {
		if !doc.IsCollection() {
			t.Error("expected IsCollection to return true")
		}
	})
	doc.Facets = []string{"NotCollectionMember"}
	t.Run("IsCollectable", func(t *testing.T) {
		if !doc.IsCollectable() {
			t.Error("expected IsCollectable to return true")
		}
	})
}

func TestPropertyMethods(t *testing.T) {
	doc := NewDocument("File", "myfile")
	// Set property
	doc.SetProperty("dc:title", NewStringField("Test Title"))
	field, found := doc.Property("dc:title")
	if !found {
		t.Error("expected property 'dc:title' to be found")
	}
	str, err := field.String()
	if err != nil || str == nil || *str != "Test Title" {
		t.Errorf("expected property value 'Test Title', got '%v', err: %v", str, err)
	}
	// Update property
	doc.SetProperty("dc:title", NewStringField("New Title"))
	field, _ = doc.Property("dc:title")
	str, err = field.String()
	if err != nil || str == nil || *str != "New Title" {
		t.Errorf("expected updated property value 'New Title', got '%v', err: %v", str, err)
	}
	// Nonexistent property
	_, found = doc.Property("dc:nonexistent")
	if found {
		t.Error("expected nonexistent property to not be found")
	}
}

func TestFileContentAndThumbnail(t *testing.T) {
	doc := NewDocument("File", "myfile")
	// Create a blob and marshal to Field
	b := blob{Filename: "file.txt", MimeType: "text/plain", Length: "123"}
	blobData, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("failed to marshal blob: %v", err)
	}
	fieldBlob := Field(blobData)
	doc.Properties[DocumentPropertyFileContent] = fieldBlob
	doc.Properties[DocumentPropertyThumbThumbnail] = fieldBlob

	t.Run("FileContent present", func(t *testing.T) {
		fc := doc.FileContent()
		if fc == nil {
			t.Fatal("expected FileContent to return blob, got nil")
		}
		if fc.Filename != "file.txt" || fc.MimeType != "text/plain" || fc.Length != "123" {
			t.Errorf("unexpected blob content: %+v", fc)
		}
	})
	t.Run("Thumbnail present", func(t *testing.T) {
		thumb := doc.Thumbnail()
		if thumb == nil {
			t.Fatal("expected Thumbnail to return blob, got nil")
		}
		if !reflect.DeepEqual(thumb, doc.FileContent()) {
			t.Errorf("expected Thumbnail to match FileContent")
		}
	})
	// Remove properties
	doc.Properties = map[string]Field{}
	t.Run("FileContent missing", func(t *testing.T) {
		if doc.FileContent() != nil {
			t.Error("expected FileContent to return nil when missing")
		}
	})
	t.Run("Thumbnail missing", func(t *testing.T) {
		if doc.Thumbnail() != nil {
			t.Error("expected Thumbnail to return nil when missing")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	doc := NewDocument("File", "myfile")
	t.Run("Empty facets", func(t *testing.T) {
		doc.Facets = nil
		if doc.HasFacet("Folderish") {
			t.Error("expected HasFacet to return false for empty facets")
		}
	})
	t.Run("Nil property value", func(t *testing.T) {
		doc.SetProperty("dc:title", NewStringField("null"))
		field, found := doc.Property("dc:title")
		if !found {
			t.Error("expected property to be found")
		}
		if !field.IsNull() {
			t.Error("expected property to be null")
		}
	})
}
