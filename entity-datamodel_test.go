package nuxeo

import "testing"

func TestEntitySchema_GetPrefix(t *testing.T) {
	cases := []struct {
		name          string
		prefix        string
		prefixAliased string
		expected      string
	}{
		{"Prefix set", "mainPrefix", "aliasPrefix", "mainPrefix"},
		{"Prefix empty, use aliased", "", "aliasPrefix", "aliasPrefix"},
		{"Both empty", "", "", ""},
	}
	for _, tc := range cases {
		s := Schema{Prefix: tc.prefix, PrefixAliased: tc.prefixAliased}
		if got := s.GetPrefix(); got != tc.expected {
			t.Errorf("%s: GetPrefix() = %q, want %q", tc.name, got, tc.expected)
		}
	}
}

func TestEntitySchemaField_TypeChecks(t *testing.T) {
	cases := []struct {
		name      string
		dataType  string
		isBlob    bool
		isBool    bool
		isComplex bool
		isDate    bool
		isLong    bool
		isDouble  bool
		isString  bool
	}{
		{"Blob type", "blob", true, false, false, false, false, false, false},
		{"Boolean type", "boolean", false, true, false, false, false, false, false},
		{"Complex type", "complex", false, false, true, false, false, false, false},
		{"Date type", "date", false, false, false, true, false, false, false},
		{"Long type", "long", false, false, false, false, true, false, false},
		{"Double type", "double", false, false, false, false, false, true, false},
		{"String type", "string", false, false, false, false, false, false, true},
		{"Unknown type", "other", false, false, false, false, false, false, false},
	}
	for _, tc := range cases {
		f := SchemaField{DataType: tc.dataType}
		if f.IsBlob() != tc.isBlob {
			t.Errorf("%s: IsBlob() = %v, want %v", tc.name, f.IsBlob(), tc.isBlob)
		}
		if f.IsBoolean() != tc.isBool {
			t.Errorf("%s: IsBoolean() = %v, want %v", tc.name, f.IsBoolean(), tc.isBool)
		}
		if f.IsComplex() != tc.isComplex {
			t.Errorf("%s: IsComplex() = %v, want %v", tc.name, f.IsComplex(), tc.isComplex)
		}
		if f.IsDate() != tc.isDate {
			t.Errorf("%s: IsDate() = %v, want %v", tc.name, f.IsDate(), tc.isDate)
		}
		if f.IsLong() != tc.isLong {
			t.Errorf("%s: IsLong() = %v, want %v", tc.name, f.IsLong(), tc.isLong)
		}
		if f.IsDouble() != tc.isDouble {
			t.Errorf("%s: IsDouble() = %v, want %v", tc.name, f.IsDouble(), tc.isDouble)
		}
		if f.IsString() != tc.isString {
			t.Errorf("%s: IsString() = %v, want %v", tc.name, f.IsString(), tc.isString)
		}
	}
}

func TestEntitySchemaField_Initialization(t *testing.T) {
	f := SchemaField{DataType: "string", IsArray: true, Fields: nil}
	if f.DataType != "string" {
		t.Errorf("DataType = %q, want %q", f.DataType, "string")
	}
	if !f.IsArray {
		t.Errorf("IsArray = %v, want true", f.IsArray)
	}
	if f.Fields != nil {
		t.Errorf("Fields = %v, want nil", f.Fields)
	}
}

func TestEntitySchema_Initialization(t *testing.T) {
	fields := map[string]SchemaField{
		"field1": {DataType: "string"},
	}
	s := Schema{Name: "schema1", Prefix: "pfx", PrefixAliased: "pfxAlias", Fields: fields}
	if s.Name != "schema1" {
		t.Errorf("Name = %q, want %q", s.Name, "schema1")
	}
	if s.Prefix != "pfx" {
		t.Errorf("Prefix = %q, want %q", s.Prefix, "pfx")
	}
	if s.PrefixAliased != "pfxAlias" {
		t.Errorf("PrefixAliased = %q, want %q", s.PrefixAliased, "pfxAlias")
	}
	if s.Fields["field1"].DataType != "string" {
		t.Errorf("Fields[\"field1\"].DataType = %q, want %q", s.Fields["field1"].DataType, "string")
	}
}

func TestEntityFacet_Initialization(t *testing.T) {
	schemas := []Schema{{Name: "schemaA"}, {Name: "schemaB"}}
	f := Facet{Name: "facet1", Schemas: schemas}
	if f.Name != "facet1" {
		t.Errorf("Name = %q, want %q", f.Name, "facet1")
	}
	if len(f.Schemas) != 2 {
		t.Errorf("Schemas length = %d, want 2", len(f.Schemas))
	}
	if f.Schemas[0].Name != "schemaA" {
		t.Errorf("Schemas[0].Name = %q, want %q", f.Schemas[0].Name, "schemaA")
	}
}

func TestEntityDocType_Initialization(t *testing.T) {
	schemas := []Schema{{Name: "schemaX"}}
	dt := DocType{Name: "docType1", Parent: "parentType", Facets: []string{"facetA"}, Schemas: schemas}
	if dt.Name != "docType1" {
		t.Errorf("Name = %q, want %q", dt.Name, "docType1")
	}
	if dt.Parent != "parentType" {
		t.Errorf("Parent = %q, want %q", dt.Parent, "parentType")
	}
	if len(dt.Facets) != 1 || dt.Facets[0] != "facetA" {
		t.Errorf("Facets = %v, want [facetA]", dt.Facets)
	}
	if len(dt.Schemas) != 1 || dt.Schemas[0].Name != "schemaX" {
		t.Errorf("Schemas[0].Name = %q, want %q", dt.Schemas[0].Name, "schemaX")
	}
}

func TestEntitySchemas_AndEntityFacets_Initialization(t *testing.T) {
	schemas := Schemas{{Name: "schema1"}, {Name: "schema2"}}
	facets := Facets{{Name: "facet1"}, {Name: "facet2"}}
	if len(schemas) != 2 {
		t.Errorf("entitySchemas length = %d, want 2", len(schemas))
	}
	if len(facets) != 2 {
		t.Errorf("entityFacets length = %d, want 2", len(facets))
	}
}

func TestEntitySchemaField_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		expect  SchemaField
		wantErr bool
	}{
		{
			name:    "String type",
			input:   `"string"`,
			expect:  SchemaField{DataType: "string", IsArray: false, Fields: nil},
			wantErr: false,
		},
		{
			name:    "String array type",
			input:   `"string[]"`,
			expect:  SchemaField{DataType: "string", IsArray: true, Fields: nil},
			wantErr: false,
		},
		{
			name:    "Complex object",
			input:   `{"type": "complex", "fields": {"sub": {"type": "string"}}}`,
			expect:  SchemaField{DataType: "complex", IsArray: false, Fields: map[string]SchemaField{"sub": {DataType: "string"}}},
			wantErr: false,
		},
		{
			name:    "Complex array object",
			input:   `{"type": "complex[]", "fields": {"sub": {"type": "string"}}}`,
			expect:  SchemaField{DataType: "complex", IsArray: true, Fields: map[string]SchemaField{"sub": {DataType: "string"}}},
			wantErr: false,
		},
		{
			name:    "Malformed JSON",
			input:   `{"type": "string"`, // missing closing brace
			expect:  SchemaField{},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var f SchemaField
			err := f.UnmarshalJSON([]byte(tc.input))
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if f.DataType != tc.expect.DataType {
				t.Errorf("DataType = %q, want %q", f.DataType, tc.expect.DataType)
			}
			if f.IsArray != tc.expect.IsArray {
				t.Errorf("IsArray = %v, want %v", f.IsArray, tc.expect.IsArray)
			}
			if len(f.Fields) != len(tc.expect.Fields) {
				t.Errorf("Fields len = %d, want %d", len(f.Fields), len(tc.expect.Fields))
			}
			for k, v := range tc.expect.Fields {
				if fv, ok := f.Fields[k]; !ok || fv.DataType != v.DataType {
					t.Errorf("Fields[%q].DataType = %q, want %q", k, fv.DataType, v.DataType)
				}
			}
		})
	}
}

func TestEntityDocTypes_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		docType string
		schema  string
		prefix  string
		wantErr bool
	}{
		{
			name:    "Valid input with docType and schema",
			input:   `{"docTypes": {"File": {"parent": "Document", "facets": ["FacetA"], "schemas": ["file"]}}, "schemas": {"file": {"@prefix": {"type": "file"}, "title": {"type": "string"}}}}`,
			docType: "File",
			schema:  "file",
			prefix:  "file",
			wantErr: false,
		},
		{
			name:    "Missing docTypes",
			input:   `{"schemas": {"file": {"@prefix": {"type": "file"}}}}`,
			docType: "",
			schema:  "file",
			prefix:  "file",
			wantErr: false,
		},
		{
			name:    "Malformed JSON",
			input:   `{"docTypes": {"File": {"parent": "Document"`, // missing closing braces
			docType: "",
			schema:  "",
			prefix:  "",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dt DocTypes
			err := dt.UnmarshalJSON([]byte(tc.input))
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.docType != "" {
				if _, ok := dt.DocTypes[tc.docType]; !ok {
					t.Errorf("DocTypes[%q] missing", tc.docType)
				}
			}
			if tc.schema != "" {
				if s, ok := dt.Schemas[tc.schema]; !ok {
					t.Errorf("Schemas[%q] missing", tc.schema)
				} else if s.Prefix != tc.prefix {
					t.Errorf("Schemas[%q].Prefix = %q, want %q", tc.schema, s.Prefix, tc.prefix)
				}
			}
		})
	}
}

func TestEntityDocTypes_Initialization(t *testing.T) {
	docTypes := map[string]DocType{"typeA": {Name: "typeA"}}
	schemas := map[string]Schema{"schemaA": {Name: "schemaA"}}
	dt := DocTypes{DocTypes: docTypes, Schemas: schemas}
	if dt.DocTypes["typeA"].Name != "typeA" {
		t.Errorf("DocTypes[\"typeA\"].Name = %q, want %q", dt.DocTypes["typeA"].Name, "typeA")
	}
	if dt.Schemas["schemaA"].Name != "schemaA" {
		t.Errorf("Schemas[\"schemaA\"].Name = %q, want %q", dt.Schemas["schemaA"].Name, "schemaA")
	}
}
