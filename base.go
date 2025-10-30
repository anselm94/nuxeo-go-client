package nuxeo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

////////////////////////////////
//// Base Entity & Entities ////
////////////////////////////////

// EntityType represents the Nuxeo entity type string (e.g., "document", "user").
type EntityType string

// entity is the base struct for all Nuxeo API entities.
// It includes the entity type and context parameters.
type entity struct {
	EntityType        EntityType       `json:"entity-type"`
	ContextParameters map[string]Field `json:"contextParameters"`
}

// ContextParameter returns the context parameter value for the given key.
func (e entity) ContextParameter(key string) (Field, bool) {
	value, found := e.ContextParameters[key]
	return value, found
}

// entities is a generic container for a list of Nuxeo entities of type T.
type entities[T any] struct {
	entity
	Entries []T `json:"entries"`
}

// paginableEntities is a generic container for paginated Nuxeo entities of type T.
type paginableEntities[T any] struct {
	entity
	IsPaginable             bool  `json:"isPaginable"`
	ResultsCount            int   `json:"resultsCount"`
	PageSize                int   `json:"pageSize"`
	MaxPageSize             int   `json:"maxPageSize"`
	ResultsCountLimit       int   `json:"resultsCountLimit"`
	CurrentPageSize         int   `json:"currentPageSize"`
	CurrentPageIndex        int   `json:"currentPageIndex"`
	CurrentPageOffset       int   `json:"currentPageOffset"`
	NumberOfPages           int   `json:"numberOfPages"`
	IsPreviousPageAvailable bool  `json:"isPreviousPageAvailable"`
	IsNextPageAvailable     bool  `json:"isNextPageAvailable"`
	IsLastPageAvailable     bool  `json:"isLastPageAvailable"`
	IsSortable              bool  `json:"isSortable"`
	HasError                bool  `json:"hasError"`
	ErrorMessage            Field `json:"errorMessage"`
	PageIndex               int   `json:"pageIndex"`
	PageCount               int   `json:"pageCount"`
	Entries                 []T   `json:"entries"`
}

////////////////////
//// Properties ////
////////////////////

// Field represents a Nuxeo property value, supporting dynamic types and JSON encoding.
type Field json.RawMessage

func NewField(value any) (Field, error) {
	if value == nil {
		return Field(json.RawMessage("null")), nil
	}

	if data, err := json.Marshal(value); err == nil {
		return Field(json.RawMessage(data)), nil
	} else {
		return Field{}, fmt.Errorf("failed to marshal field value: %w", err)
	}
}

// NewStringField creates a new Field instance representing a Nuxeo String property.
func NewStringField(value string) Field {
	val, _ := NewField(value)
	return val
}

// NewStringListField creates a new Field instance representing a Nuxeo String List property.
func NewStringListField(value []string) Field {
	val, _ := NewField(value)
	return val
}

// NewIntegerField creates a new Field instance representing a Nuxeo Long property.
func NewIntegerField(value int) Field {
	val, _ := NewField(value)
	return val
}

// NewIntegerListField creates a new Field instance representing a Nuxeo Long List property.
func NewIntegerListField(value []int) Field {
	val, _ := NewField(value)
	return val
}

// NewFloatField creates a new Field instance representing a Nuxeo Double property.
func NewFloatField(value float64) Field {
	val, _ := NewField(value)
	return val
}

// NewFloatListField creates a new Field instance representing a Nuxeo Double List property.
func NewFloatListField(value []float64) Field {
	val, _ := NewField(value)
	return val
}

// NewBooleanField creates a new Field instance representing a Nuxeo Boolean property.
func NewBooleanField(value bool) Field {
	val, _ := NewField(value)
	return val
}

// NewBooleanListField creates a new Field instance representing a Nuxeo Boolean List property.
func NewBooleanListField(value []bool) Field {
	val, _ := NewField(value)
	return val
}

// NewTimeField creates a new Field instance representing a Nuxeo Calendar property.
func NewTimeField(value ISO8601Time) Field {
	val, _ := NewField(value)
	return val
}

// NewTimeListField creates a new Field instance representing a Nuxeo Calendar List property.
func NewTimeListField(value []ISO8601Time) Field {
	val, _ := NewField(value)
	return val
}

// NewComplexField creates a new Field instance representing a Nuxeo Complex Property.
func NewComplexField(value any) (Field, error) {
	return NewField(value)
}

// NewComplexListField creates a new Field instance representing a Nuxeo Complex Property List.
func NewComplexListField(value []any) (Field, error) {
	return NewField(value)
}

// MarshalJSON implements the json.Marshaler interface for Field.
func (f Field) MarshalJSON() ([]byte, error) {
	return json.RawMessage(f).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface for Field.
func (f *Field) UnmarshalJSON(data []byte) error {
	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*f = Field(raw)
	return nil
}

// IsNull returns true if the Field is null.
func (f Field) IsNull() bool {
	return string(f) == "\"null\"" || string(f) == "null"
}

// String returns the Field value as a string (Nuxeo String).
func (f Field) String() (*string, error) {
	if f.IsNull() {
		return nil, nil
	}
	var s string
	if err := json.Unmarshal(f, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// StringList returns the Field value as a slice of strings.
func (f Field) StringList() ([]string, error) {
	if f.IsNull() {
		return nil, nil
	}

	var stringList []string
	if err := json.Unmarshal(f, &stringList); err != nil {
		return nil, err
	}
	return stringList, nil
}

// Integer returns the Field value as an integer (Nuxeo Long).
func (f Field) Integer() (*int, error) {
	if f.IsNull() {
		return nil, nil
	}

	var intValue int
	if err := json.Unmarshal(f, &intValue); err != nil {
		return nil, err
	}
	return &intValue, nil
}

// IntegerList returns the Field value as a slice of integers.
func (f Field) IntegerList() ([]int, error) {
	if f.IsNull() {
		return nil, nil
	}

	var intList []int
	if err := json.Unmarshal(f, &intList); err != nil {
		return nil, err
	}
	return intList, nil
}

// Float returns the Field value as a float64 (Nuxeo Double).
func (f Field) Float() (*float64, error) {
	if f.IsNull() {
		return nil, nil
	}

	var floatValue float64
	if err := json.Unmarshal(f, &floatValue); err != nil {
		return nil, err
	}
	return &floatValue, nil
}

// FloatList returns the Field value as a slice of float64.
func (f Field) FloatList() ([]float64, error) {
	if f.IsNull() {
		return nil, nil
	}

	var floatList []float64
	if err := json.Unmarshal(f, &floatList); err != nil {
		return nil, err
	}
	return floatList, nil
}

// Boolean returns the Field value as a bool (Nuxeo Boolean).
func (f Field) Boolean() (*bool, error) {
	if f.IsNull() {
		return nil, nil
	}

	var boolValue bool
	if err := json.Unmarshal(f, &boolValue); err != nil {
		return nil, err
	}
	return &boolValue, nil
}

// BooleanList returns the Field value as a slice of bool.
func (f Field) BooleanList() ([]bool, error) {
	if f.IsNull() {
		return nil, nil
	}

	var boolList []bool
	if err := json.Unmarshal(f, &boolList); err != nil {
		return nil, err
	}
	return boolList, nil
}

// Time returns the Field value as an ISO8601Time (Nuxeo Calendar).
func (f Field) Time() (*ISO8601Time, error) {
	if f.IsNull() {
		return nil, nil
	}

	var timeValue ISO8601Time
	if err := json.Unmarshal(f, &timeValue); err != nil {
		return nil, err
	}
	return &timeValue, nil
}

// TimeList returns the Field value as a slice of ISO8601Time.
func (f Field) TimeList() ([]ISO8601Time, error) {
	if f.IsNull() {
		return nil, nil
	}

	var timeList []ISO8601Time
	if err := json.Unmarshal(f, &timeList); err != nil {
		return nil, err
	}
	return timeList, nil
}

// Complex decodes the Field into the provided output struct (Nuxeo Complex Property).
func (f Field) Complex(out any) error {
	if f.IsNull() {
		return nil
	}

	return json.Unmarshal(f, out)
}

// ComplexList decodes the Field into the provided output slice (Nuxeo Complex Property List).
func (f Field) ComplexList(out any) error {
	if f.IsNull() {
		return nil
	}

	return json.Unmarshal(f, out)
}

//////////////////////////////
//// Pagination & Sorting ////
//////////////////////////////

// PaginationOptions specifies pagination parameters for Nuxeo API requests.
type PaginationOptions struct {
	CurrentPageIndex int `json:"currentPageIndex"`
	PageSize         int `json:"pageSize"`
}

// QueryParams returns the pagination options as URL query parameters.
func (p *PaginationOptions) QueryParams() url.Values {
	if p == nil {
		return nil
	}
	queryParams := make(url.Values)
	if p.CurrentPageIndex > -1 {
		queryParams.Set("currentPageIndex", fmt.Sprintf("%d", p.CurrentPageIndex))
	}
	if p.PageSize != 0 {
		queryParams.Set("pageSize", fmt.Sprintf("%d", p.PageSize))
	}
	return queryParams
}

// SortedPaginationOptions specifies pagination and sorting parameters for Nuxeo API requests.
type SortedPaginationOptions struct {
	CurrentPageIndex int    `json:"currentPageIndex"`
	PageSize         int    `json:"pageSize"`
	MaxResults       int    `json:"maxResults"`
	SortBy           string `json:"sortBy"`
	SortOrder        string `json:"sortOrder"`
}

// QueryParams returns the sorted pagination options as URL query parameters.
func (p *SortedPaginationOptions) QueryParams() url.Values {
	if p == nil {
		return nil
	}
	queryParams := make(url.Values)
	if p.CurrentPageIndex > -1 {
		queryParams.Set("currentPageIndex", fmt.Sprintf("%d", p.CurrentPageIndex))
	}
	if p.PageSize != 0 {
		queryParams.Set("pageSize", fmt.Sprintf("%d", p.PageSize))
	}
	if p.MaxResults != 0 {
		queryParams.Set("maxResults", fmt.Sprintf("%d", p.MaxResults))
	}
	if p.SortBy != "" {
		queryParams.Set("sortBy", p.SortBy)
	}
	if p.SortOrder != "" {
		queryParams.Set("sortOrder", p.SortOrder)
	}
	return queryParams
}

//////////////////////
//// ISO8601 Time ////
//////////////////////

// ISO8601Time represents a time in ISO8601 format (Nuxeo Calendar property).
type ISO8601Time time.Time

// ISO8601TimeLayout is the layout string for ISO8601 time formatting.
const ISO8601TimeLayout = "2006-01-02T15:04:05.999Z"

// UnmarshalJSON parses an ISO8601Time from JSON.
func (t *ISO8601Time) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse(`"`+ISO8601TimeLayout+`"`, string(data))
	if err != nil {
		return err
	}
	*t = ISO8601Time(parsedTime)
	return nil
}

// MarshalJSON formats an ISO8601Time as JSON.
func (t ISO8601Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(ISO8601TimeLayout) + `"`), nil
}
