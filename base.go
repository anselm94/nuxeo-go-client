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

type EntityType string

type entity struct {
	EntityType        EntityType       `json:"entity-type"`
	ContextParameters map[string]Field `json:"contextParameters"`
}

func (e entity) ContextParameter(key string) Field {
	return e.ContextParameters[key]
}

type entities[T any] struct {
	entity
	Entries []T `json:"entries"`
}

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

func (f Field) MarshalJSON() ([]byte, error) {
	return json.RawMessage(f).MarshalJSON()
}

func (f *Field) UnmarshalJSON(data []byte) error {
	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*f = Field(raw)
	return nil
}

func (f Field) IsNull() bool {
	return string(f) == "null"
}

// String returns the Field value as a string - Nuxeo String
func (f Field) String() (*string, error) {
	if f.IsNull() {
		return nil, nil
	}

	strValue := string(f)
	return &strValue, nil
}

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

// Integer returns the Field value as an integer - Nuxeo Long
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

// Float returns the Field value as a float - Nuxeo Double
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

// Boolean returns the Field value as a boolean - Nuxeo Boolean
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

// Time returns the Field value as an ISO8601Time - Nuxeo Calendar
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

// Complex decodes the Field into the provided output struct - Nuxeo Complex Property
func (f Field) Complex(out any) error {
	if f.IsNull() {
		out = nil // Set output to nil if field is null
		return nil
	}

	return json.Unmarshal(f, out)
}

func (f Field) ComplexList(out any) error {
	if f.IsNull() {
		out = nil // Set output to nil if field is null
		return nil
	}

	return json.Unmarshal(f, out)
}

//////////////////////////////
//// Pagination & Sorting ////
//////////////////////////////

type PaginationOptions struct {
	CurrentPageIndex int `json:"currentPageIndex"`
	PageSize         int `json:"pageSize"`
}

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

type SortedPaginationOptions struct {
	CurrentPageIndex int    `json:"currentPageIndex"`
	PageSize         int    `json:"pageSize"`
	MaxResults       int    `json:"maxResults"`
	SortBy           string `json:"sortBy"`
	SortOrder        string `json:"sortOrder"`
}

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

type ISO8601Time time.Time

const ISO8601TimeLayout = "2006-01-02T15:04:05.999Z"

func (t *ISO8601Time) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse(`"`+ISO8601TimeLayout+`"`, string(data))
	if err != nil {
		return err
	}
	*t = ISO8601Time(parsedTime)
	return nil
}

func (t ISO8601Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(ISO8601TimeLayout) + `"`), nil
}
