package nuxeo

import (
	"encoding/json"
	"strings"
)

type Schema struct {
	EntityType    string                 `json:"entity-type,omitempty"`
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

func (sf SchemaField) IsString() bool {
	return sf.DataType == "string"
}

func (sf *SchemaField) UnmarshalJSON(data []byte) error {
	// if value is a string
	if n := len(data); data[0] == '"' && data[n-1] == '"' {
		strVal := strings.Trim(string(data), "\"")
		isMultiValued := strings.HasSuffix(strVal, "[]")
		strVal = strings.TrimSuffix(strVal, "[]")
		*sf = SchemaField{
			DataType: strVal,
			IsArray:  isMultiValued,
			Fields:   nil,
		}
	}

	// if value is an object
	var complexSchemaField struct {
		DataType string                 `json:"type"`
		Fields   map[string]SchemaField `json:"fields"`
	}
	if err := json.Unmarshal(data, &complexSchemaField); err == nil {
		isMultiValued := strings.HasSuffix(complexSchemaField.DataType, "[]")
		complexSchemaField.DataType = strings.TrimSuffix(complexSchemaField.DataType, "[]")
		*sf = SchemaField{
			DataType: complexSchemaField.DataType,
			IsArray:  isMultiValued,
			Fields:   complexSchemaField.Fields,
		}
	}
	return nil
}

type Schemas []Schema

type Facet struct {
	EntityType string   `json:"entity-type,omitempty"`
	Name       string   `json:"name,omitempty"`
	Schemas    []Schema `json:"schemas"`
}

type Facets []Facet

type DocType struct {
	EntityType string   `json:"entity-type,omitempty"`
	Name       string   `json:"name,omitempty"`
	Parent     string   `json:"parent"`
	Facets     []string `json:"facets"`
	Schemas    []Schema `json:"schemas"` // can be []string or "Schema" based on get all doc types vs get single doc type
}

type DocTypes struct {
	DocTypes map[string]DocType
	Schemas  map[string]Schema
}

func (dt *DocTypes) UnmarshalJSON(data []byte) error {
	var docTypes struct {
		DocTypes map[string]struct {
			Parent  string   `json:"parent"`
			Facets  []string `json:"facets"`
			Schemas []string `json:"schemas"`
		} `json:"docTypes"`
		Schemas map[string]map[string]SchemaField `json:"schemas"`
	}
	if err := json.Unmarshal(data, &docTypes); err != nil {
		return err
	}

	*dt = DocTypes{
		DocTypes: make(map[string]DocType, len(docTypes.DocTypes)),
		Schemas:  make(map[string]Schema, len(docTypes.Schemas)),
	}

	for name, schemaFields := range docTypes.Schemas {
		prefixField := schemaFields["@prefix"]
		delete(schemaFields, "@prefix") // "@prefix" is not an actual field of the schema but Nuxeo uses it to convey the prefix info
		(*dt).Schemas[name] = Schema{
			Name:          name,
			Prefix:        prefixField.DataType, // "@prefix" field holds the prefix string in its DataType
			PrefixAliased: prefixField.DataType, // "@prefix" field holds the prefix string in its DataType
			Fields:        schemaFields,
		}
	}

	for name, docType := range docTypes.DocTypes {
		docTypeSchemas := make([]Schema, 0, len(docType.Schemas))
		for _, schemaName := range docType.Schemas {
			if schema, exists := (*dt).Schemas[schemaName]; exists {
				docTypeSchemas = append(docTypeSchemas, schema)
			}
		}
		(*dt).DocTypes[name] = DocType{
			Name:    name,
			Parent:  docType.Parent,
			Facets:  docType.Facets,
			Schemas: docTypeSchemas,
		}
	}

	return nil
}
