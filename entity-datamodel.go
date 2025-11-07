package nuxeo

import (
	"encoding/json"
	"strings"
)

// Schema represents a Nuxeo schema entity as returned by the Data Model API.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-schema
// Fields:
//   - Name: schema name
//   - Prefix: schema prefix (for property names)
//   - PrefixAliased: alternative prefix field (from Nuxeo API)
//   - Fields: map of field names to schema field definitions
type Schema struct {
	entity
	Name          string                 `json:"name"`
	Prefix        string                 `json:"prefix"`
	PrefixAliased string                 `json:"@prefix"`
	Fields        map[string]SchemaField `json:"fields"`
}

func (s Schema) GetPrefix() string {
	if s.Prefix != "" {
		return s.Prefix
	}
	return s.PrefixAliased
}

// SchemaField describes a field in a Nuxeo schema.
// DataType: type of the field (e.g. string, boolean, date, blob, complex)
// IsArray: true if the field is an array
// Fields: for complex types, nested fields
type SchemaField struct {
	DataType string
	IsArray  bool
	Fields   map[string]SchemaField
}

func (sf SchemaField) IsBlob() bool {
	return sf.DataType == "blob"
}

func (sf SchemaField) IsBoolean() bool {
	return sf.DataType == "boolean"
}

func (sf SchemaField) IsComplex() bool {
	return sf.DataType == "complex"
}

func (sf SchemaField) IsDate() bool {
	return sf.DataType == "date"
}

func (sf SchemaField) IsLong() bool {
	return sf.DataType == "long"
}

func (sf SchemaField) IsDouble() bool {
	return sf.DataType == "double"
}

func (sf SchemaField) IsString() bool {
	return sf.DataType == "string"
}

func (sf *SchemaField) UnmarshalJSON(data []byte) error {
	strLength := len(data)
	if strLength == 0 {
		return json.Unmarshal(data, sf) // let json handle empty input
	}

	// if value is a string
	if data[0] == '"' && data[strLength-1] == '"' {
		strVal := strings.Trim(string(data), "\"")
		isMultiValued := strings.HasSuffix(strVal, "[]")
		strVal = strings.TrimSuffix(strVal, "[]")
		*sf = SchemaField{
			DataType: strVal,
			IsArray:  isMultiValued,
			Fields:   nil,
		}
		return nil
	}

	// if value is a complex object
	var complexSchemaField struct {
		DataType string                 `json:"type"`
		Fields   map[string]SchemaField `json:"fields"`
	}
	if err := json.Unmarshal(data, &complexSchemaField); err != nil {
		return err
	}

	isMultiValued := strings.HasSuffix(complexSchemaField.DataType, "[]")
	complexSchemaField.DataType = strings.TrimSuffix(complexSchemaField.DataType, "[]")
	*sf = SchemaField{
		DataType: complexSchemaField.DataType,
		IsArray:  isMultiValued,
		Fields:   complexSchemaField.Fields,
	}
	return nil
}

// Schemas is a slice of entitySchema, representing a collection of schemas.
type Schemas []Schema

// Facet represents a Nuxeo facet entity as returned by the Data Model API.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-facet
// Fields:
//   - Name: facet name
//   - Schemas: schemas associated with the facet
type Facet struct {
	entity
	Name    string   `json:"name,omitempty"`
	Schemas []Schema `json:"schemas"`
}

// Facets is a slice of entityFacet, representing a collection of facets.
type Facets []Facet

// DocType represents a Nuxeo document type entity as returned by the Data Model API.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-document-type
// Fields:
//   - Name: document type name
//   - Parent: parent document type
//   - Facets: list of facet names
//   - Schemas: schemas associated with the document type
type DocType struct {
	entity
	Name    string   `json:"name,omitempty"`
	Parent  string   `json:"parent"`
	Facets  []string `json:"facets"`
	Schemas []Schema `json:"schemas"` // can be []string or "Schema" based on get all doc types vs get single doc type
}

// DocTypes represents the response from GET /config/types.
// Contains all document types and their associated schemas.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-all-document-types
type DocTypes struct {
	DocTypes map[string]DocType `json:"docTypes"`
	Schemas  map[string]Schema  `json:"schemas"`
}

func (dt *DocTypes) UnmarshalJSON(data []byte) error {
	var vDocTypes struct {
		DocTypes map[string]struct {
			Parent  string   `json:"parent"`
			Facets  []string `json:"facets"`
			Schemas []string `json:"schemas"`
		} `json:"docTypes"`
		Schemas map[string]map[string]SchemaField `json:"schemas"`
	}
	if err := json.Unmarshal(data, &vDocTypes); err != nil {
		return err
	}

	*dt = DocTypes{
		DocTypes: make(map[string]DocType, len(vDocTypes.DocTypes)),
		Schemas:  make(map[string]Schema, len(vDocTypes.Schemas)),
	}

	for name, schemaFields := range vDocTypes.Schemas {
		prefixField := schemaFields["@prefix"]
		delete(schemaFields, "@prefix") // "@prefix" is not an actual field of the schema but Nuxeo uses it to convey the prefix info
		(*dt).Schemas[name] = Schema{
			Name:          name,
			Prefix:        prefixField.DataType, // "@prefix" field holds the prefix string in its DataType
			PrefixAliased: prefixField.DataType, // "@prefix" field holds the prefix string in its DataType
			Fields:        schemaFields,
		}
	}

	for name, vDocType := range vDocTypes.DocTypes {
		docTypeSchemas := make([]Schema, 0, len(vDocType.Schemas))
		for _, schemaName := range vDocType.Schemas {
			if schema, exists := (*dt).Schemas[schemaName]; exists {
				docTypeSchemas = append(docTypeSchemas, schema)
			}
		}
		(*dt).DocTypes[name] = DocType{
			Name:    name,
			Parent:  vDocType.Parent,
			Facets:  vDocType.Facets,
			Schemas: docTypeSchemas,
		}
	}

	return nil
}
