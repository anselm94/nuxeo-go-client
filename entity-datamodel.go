package nuxeo

import (
	"encoding/json"
	"strings"
)

// entitySchema represents a Nuxeo schema entity as returned by the Data Model API.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-schema
// Fields:
//   - Name: schema name
//   - Prefix: schema prefix (for property names)
//   - PrefixAliased: alternative prefix field (from Nuxeo API)
//   - Fields: map of field names to schema field definitions
type entitySchema struct {
	entity
	Name          string                       `json:"name"`
	Prefix        string                       `json:"prefix"`
	PrefixAliased string                       `json:"@prefix"`
	Fields        map[string]entitySchemaField `json:"fields"`
}

func (s entitySchema) GetPrefix() string {
	if s.Prefix != "" {
		return s.Prefix
	}
	return s.PrefixAliased
}

// entitySchemaField describes a field in a Nuxeo schema.
// DataType: type of the field (e.g. string, boolean, date, blob, complex)
// IsArray: true if the field is an array
// Fields: for complex types, nested fields
type entitySchemaField struct {
	DataType string
	IsArray  bool
	Fields   map[string]entitySchemaField
}

func (sf entitySchemaField) IsBlob() bool {
	return sf.DataType == "blob"
}

func (sf entitySchemaField) IsBoolean() bool {
	return sf.DataType == "boolean"
}

func (sf entitySchemaField) IsComplex() bool {
	return sf.DataType == "complex"
}

func (sf entitySchemaField) IsDate() bool {
	return sf.DataType == "date"
}

func (sf entitySchemaField) IsLong() bool {
	return sf.DataType == "long"
}

func (sf entitySchemaField) IsString() bool {
	return sf.DataType == "string"
}

func (sf *entitySchemaField) UnmarshalJSON(data []byte) error {
	strLength := len(data)
	if strLength == 0 {
		return json.Unmarshal(data, sf) // let json handle empty input
	}

	// if value is a string
	if data[0] == '"' && data[strLength-1] == '"' {
		strVal := strings.Trim(string(data), "\"")
		isMultiValued := strings.HasSuffix(strVal, "[]")
		strVal = strings.TrimSuffix(strVal, "[]")
		*sf = entitySchemaField{
			DataType: strVal,
			IsArray:  isMultiValued,
			Fields:   nil,
		}
		return nil
	}

	// if value is a complex object
	var complexSchemaField struct {
		DataType string                       `json:"type"`
		Fields   map[string]entitySchemaField `json:"fields"`
	}
	if err := json.Unmarshal(data, &complexSchemaField); err != nil {
		return err
	}

	isMultiValued := strings.HasSuffix(complexSchemaField.DataType, "[]")
	complexSchemaField.DataType = strings.TrimSuffix(complexSchemaField.DataType, "[]")
	*sf = entitySchemaField{
		DataType: complexSchemaField.DataType,
		IsArray:  isMultiValued,
		Fields:   complexSchemaField.Fields,
	}
	return nil
}

// entitySchemas is a slice of entitySchema, representing a collection of schemas.
type entitySchemas []entitySchema

// entityFacet represents a Nuxeo facet entity as returned by the Data Model API.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-facet
// Fields:
//   - Name: facet name
//   - Schemas: schemas associated with the facet
type entityFacet struct {
	entity
	Name    string         `json:"name,omitempty"`
	Schemas []entitySchema `json:"schemas"`
}

// entityFacets is a slice of entityFacet, representing a collection of facets.
type entityFacets []entityFacet

// entityDocType represents a Nuxeo document type entity as returned by the Data Model API.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-document-type
// Fields:
//   - Name: document type name
//   - Parent: parent document type
//   - Facets: list of facet names
//   - Schemas: schemas associated with the document type
type entityDocType struct {
	entity
	Name    string         `json:"name,omitempty"`
	Parent  string         `json:"parent"`
	Facets  []string       `json:"facets"`
	Schemas []entitySchema `json:"schemas"` // can be []string or "Schema" based on get all doc types vs get single doc type
}

// entityDocTypes represents the response from GET /config/types.
// Contains all document types and their associated schemas.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-all-document-types
type entityDocTypes struct {
	DocTypes map[string]entityDocType `json:"docTypes"`
	Schemas  map[string]entitySchema  `json:"schemas"`
}

func (dt *entityDocTypes) UnmarshalJSON(data []byte) error {
	var vDocTypes struct {
		DocTypes map[string]struct {
			Parent  string   `json:"parent"`
			Facets  []string `json:"facets"`
			Schemas []string `json:"schemas"`
		} `json:"docTypes"`
		Schemas map[string]map[string]entitySchemaField `json:"schemas"`
	}
	if err := json.Unmarshal(data, &vDocTypes); err != nil {
		return err
	}

	*dt = entityDocTypes{
		DocTypes: make(map[string]entityDocType, len(vDocTypes.DocTypes)),
		Schemas:  make(map[string]entitySchema, len(vDocTypes.Schemas)),
	}

	for name, schemaFields := range vDocTypes.Schemas {
		prefixField := schemaFields["@prefix"]
		delete(schemaFields, "@prefix") // "@prefix" is not an actual field of the schema but Nuxeo uses it to convey the prefix info
		(*dt).Schemas[name] = entitySchema{
			Name:          name,
			Prefix:        prefixField.DataType, // "@prefix" field holds the prefix string in its DataType
			PrefixAliased: prefixField.DataType, // "@prefix" field holds the prefix string in its DataType
			Fields:        schemaFields,
		}
	}

	for name, vDocType := range vDocTypes.DocTypes {
		docTypeSchemas := make([]entitySchema, 0, len(vDocType.Schemas))
		for _, schemaName := range vDocType.Schemas {
			if schema, exists := (*dt).Schemas[schemaName]; exists {
				docTypeSchemas = append(docTypeSchemas, schema)
			}
		}
		(*dt).DocTypes[name] = entityDocType{
			Name:    name,
			Parent:  vDocType.Parent,
			Facets:  vDocType.Facets,
			Schemas: docTypeSchemas,
		}
	}

	return nil
}
