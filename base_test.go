package nuxeo

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNewField_IsNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  any
		want   string
		isNull bool
	}{
		{"nil value", nil, "null", true},
		{"string value", "foo", `"foo"`, false},
		{"int value", 42, "42", false},
		{"slice value", []string{"a", "b"}, `["a","b"]`, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewField(tc.input)
			if err != nil {
				t.Fatalf("NewField error: %v", err)
			}
			if string(f) != tc.want {
				t.Errorf("got %q, want %q", string(f), tc.want)
			}
			if got := f.IsNull(); got != tc.isNull {
				t.Errorf("IsNull() = %v, want %v", got, tc.isNull)
			}
		})
	}
}

func TestField_StringAndStringList(t *testing.T) {
	t.Parallel()

	t.Run("String: null", func(t *testing.T) {
		f, _ := NewField(nil)
		str, err := f.String()
		if err != nil {
			t.Fatalf("String() error: %v", err)
		}
		if str != nil {
			t.Errorf("String() got %v, want nil", str)
		}
	})

	t.Run("String: valid string", func(t *testing.T) {
		f := NewStringField("hello")
		str, err := f.String()
		if err != nil {
			t.Fatalf("String() error: %v", err)
		}
		if str == nil || *str != "hello" {
			t.Errorf("String() got %v, want %q", str, "hello")
		}
	})

	t.Run("StringList: null", func(t *testing.T) {
		f, _ := NewField(nil)
		lst, err := f.StringList()
		if err != nil {
			t.Fatalf("StringList() error: %v", err)
		}
		if lst != nil {
			t.Errorf("StringList() got %v, want nil", lst)
		}
	})

	t.Run("StringList: valid list", func(t *testing.T) {
		f := NewStringListField([]string{"a", "b"})
		lst, err := f.StringList()
		if err != nil {
			t.Fatalf("StringList() error: %v", err)
		}
		if len(lst) != 2 || lst[0] != "a" || lst[1] != "b" {
			t.Errorf("StringList() got %v, want [a b]", lst)
		}
	})

	// Direct usage of encoding/json to satisfy import
	var dummy []string
	_ = json.Unmarshal([]byte(`["x","y"]`), &dummy)
}

func TestField_MarshalUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		wantRaw string
		wantErr bool
	}{
		{"null", nil, "null", false},
		{"string", "bar", `"bar"`, false},
		{"int", 99, "99", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewField(tc.input)
			if err != nil {
				t.Fatalf("NewField error: %v", err)
			}
			b, err := f.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON error: %v", err)
			}
			if string(b) != tc.wantRaw {
				t.Errorf("MarshalJSON got %q, want %q", string(b), tc.wantRaw)
			}

			var f2 Field
			err = f2.UnmarshalJSON(b)
			if (err != nil) != tc.wantErr {
				t.Errorf("UnmarshalJSON error = %v, wantErr %v", err, tc.wantErr)
			}
			if string(f2) != tc.wantRaw {
				t.Errorf("UnmarshalJSON got %q, want %q", string(f2), tc.wantRaw)
			}
		})
	}

	// Invalid JSON
	t.Run("invalid json", func(t *testing.T) {
		var f Field
		err := f.UnmarshalJSON([]byte("{bad json}"))
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})
}

func TestField_NumericAndBoolMethods(t *testing.T) {
	t.Parallel()

	t.Run("Integer: null", func(t *testing.T) {
		f, _ := NewField(nil)
		val, err := f.Integer()
		if err != nil {
			t.Fatalf("Integer() error: %v", err)
		}
		if val != nil {
			t.Errorf("Integer() got %v, want nil", val)
		}
	})

	t.Run("Integer: valid", func(t *testing.T) {
		f := NewIntegerField(123)
		val, err := f.Integer()
		if err != nil {
			t.Fatalf("Integer() error: %v", err)
		}
		if val == nil || *val != 123 {
			t.Errorf("Integer() got %v, want 123", val)
		}
	})

	t.Run("IntegerList: valid", func(t *testing.T) {
		f := NewIntegerListField([]int{1, 2, 3})
		lst, err := f.IntegerList()
		if err != nil {
			t.Fatalf("IntegerList() error: %v", err)
		}
		if !reflect.DeepEqual(lst, []int{1, 2, 3}) {
			t.Errorf("IntegerList() got %v, want [1 2 3]", lst)
		}
	})

	t.Run("Float: valid", func(t *testing.T) {
		f := NewFloatField(3.14)
		val, err := f.Float()
		if err != nil {
			t.Fatalf("Float() error: %v", err)
		}
		if val == nil || *val != 3.14 {
			t.Errorf("Float() got %v, want 3.14", val)
		}
	})

	t.Run("FloatList: valid", func(t *testing.T) {
		f := NewFloatListField([]float64{1.1, 2.2})
		lst, err := f.FloatList()
		if err != nil {
			t.Fatalf("FloatList() error: %v", err)
		}
		if !reflect.DeepEqual(lst, []float64{1.1, 2.2}) {
			t.Errorf("FloatList() got %v, want [1.1 2.2]", lst)
		}
	})

	t.Run("Boolean: valid", func(t *testing.T) {
		f := NewBooleanField(true)
		val, err := f.Boolean()
		if err != nil {
			t.Fatalf("Boolean() error: %v", err)
		}
		if val == nil || *val != true {
			t.Errorf("Boolean() got %v, want true", val)
		}
	})

	t.Run("BooleanList: valid", func(t *testing.T) {
		f := NewBooleanListField([]bool{true, false})
		lst, err := f.BooleanList()
		if err != nil {
			t.Fatalf("BooleanList() error: %v", err)
		}
		if !reflect.DeepEqual(lst, []bool{true, false}) {
			t.Errorf("BooleanList() got %v, want [true false]", lst)
		}
	})
}

func TestField_TimeAndTimeList(t *testing.T) {
	t.Parallel()

	t.Run("Time: null", func(t *testing.T) {
		f, _ := NewField(nil)
		val, err := f.Time()
		if err != nil {
			t.Fatalf("Time() error: %v", err)
		}
		if val != nil {
			t.Errorf("Time() got %v, want nil", val)
		}
	})

	t.Run("Time: valid", func(t *testing.T) {
		tm := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)
		iso := ISO8601Time(tm)
		f := NewTimeField(iso)
		val, err := f.Time()
		if err != nil {
			t.Fatalf("Time() error: %v", err)
		}
		if val == nil || !time.Time(*val).Equal(tm) {
			t.Errorf("Time() got %v, want %v", val, tm)
		}
	})

	t.Run("TimeList: valid", func(t *testing.T) {
		tm1 := ISO8601Time(time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC))
		tm2 := ISO8601Time(time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC))
		f := NewTimeListField([]ISO8601Time{tm1, tm2})
		lst, err := f.TimeList()
		if err != nil {
			t.Fatalf("TimeList() error: %v", err)
		}
		if len(lst) != 2 || !time.Time(lst[0]).Equal(time.Time(tm1)) || !time.Time(lst[1]).Equal(time.Time(tm2)) {
			t.Errorf("TimeList() got %v, want [%v %v]", lst, tm1, tm2)
		}
	})
}

func TestField_ComplexAndComplexList(t *testing.T) {
	t.Parallel()

	t.Run("Complex: valid struct", func(t *testing.T) {
		type Foo struct {
			Bar string `json:"bar"`
		}
		f, _ := NewComplexField(Foo{Bar: "baz"})
		var out Foo
		err := f.Complex(&out)
		if err != nil {
			t.Fatalf("Complex() error: %v", err)
		}
		if out.Bar != "baz" {
			t.Errorf("Complex() got %v, want Bar=baz", out)
		}
	})

	t.Run("ComplexList: valid slice", func(t *testing.T) {
		type Foo struct {
			Bar string `json:"bar"`
		}
		f, _ := NewComplexField([]Foo{{Bar: "a"}, {Bar: "b"}})
		var out []Foo
		err := f.ComplexList(&out)
		if err != nil {
			t.Fatalf("ComplexList() error: %v", err)
		}
		if len(out) != 2 || out[0].Bar != "a" || out[1].Bar != "b" {
			t.Errorf("ComplexList() got %v, want [a b]", out)
		}
	})

	t.Run("Complex: null", func(t *testing.T) {
		type Foo struct{ Bar string }
		f, _ := NewComplexField(nil)
		var out Foo
		err := f.Complex(&out)
		if err != nil {
			t.Fatalf("Complex() error: %v", err)
		}
		// Output should be zero value
		if out != (Foo{}) {
			t.Errorf("Complex() got %v, want zero value", out)
		}
	})
}

func TestEntity_ContextParameter(t *testing.T) {
	t.Parallel()
	e := entity{
		EntityType: "document",
		ContextParameters: map[string]Field{
			"foo": Field([]byte(`"bar"`)),
			"num": Field([]byte(`42`)),
		},
	}

	f, _ := e.ContextParameter("foo")
	str, err := f.String()
	if err != nil {
		t.Fatalf("ContextParameter String() error: %v", err)
	}
	if str == nil || *str != "bar" {
		t.Errorf("ContextParameter got %v, want \"bar\"", str)
	}

	num, _ := e.ContextParameter("num")
	ival, err := num.Integer()
	if err != nil {
		t.Fatalf("ContextParameter Integer() error: %v", err)
	}
	if ival == nil || *ival != 42 {
		t.Errorf("ContextParameter got %v, want 42", ival)
	}

	_, ok := e.ContextParameter("missing")
	if ok {
		t.Errorf("ContextParameter missing key should be zero value Field (empty)")
	}
}

func TestPaginationOptions_QueryParams(t *testing.T) {
	t.Parallel()

	t.Run("nil options", func(t *testing.T) {
		var p *PaginationOptions
		if vals := p.QueryParams(); vals != nil {
			t.Errorf("QueryParams(nil) got %v, want nil", vals)
		}
	})

	t.Run("default values", func(t *testing.T) {
		p := &PaginationOptions{}
		vals := p.QueryParams()
		if vals.Get("currentPageIndex") != "0" || len(vals) != 1 {
			t.Errorf("QueryParams default got %v, want currentPageIndex=0 only", vals)
		}
	})

	t.Run("set values", func(t *testing.T) {
		p := &PaginationOptions{CurrentPageIndex: 2, PageSize: 50}
		vals := p.QueryParams()
		if vals.Get("currentPageIndex") != "2" || vals.Get("pageSize") != "50" {
			t.Errorf("QueryParams got %v, want currentPageIndex=2, pageSize=50", vals)
		}
	})
}

func TestSortedPaginationOptions_QueryParams(t *testing.T) {
	t.Parallel()

	t.Run("nil options", func(t *testing.T) {
		var p *SortedPaginationOptions
		if vals := p.QueryParams(); vals != nil {
			t.Errorf("QueryParams(nil) got %v, want nil", vals)
		}
	})

	t.Run("default values", func(t *testing.T) {
		p := &SortedPaginationOptions{}
		vals := p.QueryParams()
		if vals.Get("currentPageIndex") != "0" || len(vals) != 1 {
			t.Errorf("QueryParams default got %v, want currentPageIndex=0 only", vals)
		}
	})

	t.Run("set values", func(t *testing.T) {
		p := &SortedPaginationOptions{CurrentPageIndex: 1, PageSize: 25, MaxResults: 100, SortBy: "title", SortOrder: "desc"}
		vals := p.QueryParams()
		if vals.Get("currentPageIndex") != "1" || vals.Get("pageSize") != "25" || vals.Get("maxResults") != "100" || vals.Get("sortBy") != "title" || vals.Get("sortOrder") != "desc" {
			t.Errorf("QueryParams got %v, want all set", vals)
		}
	})
}

func TestISO8601Time_MarshalUnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("round trip", func(t *testing.T) {
		tm := time.Date(2025, 10, 27, 15, 4, 5, 0, time.UTC)
		iso := ISO8601Time(tm)
		b, err := iso.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON error: %v", err)
		}
		var iso2 ISO8601Time
		err = iso2.UnmarshalJSON(b)
		if err != nil {
			t.Fatalf("UnmarshalJSON error: %v", err)
		}
		if !time.Time(iso2).Equal(tm) {
			t.Errorf("UnmarshalJSON got %v, want %v", iso2, tm)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		var iso ISO8601Time
		err := iso.UnmarshalJSON([]byte(`"not-a-date"`))
		if err == nil {
			t.Error("expected error for invalid date format, got nil")
		}
	})
}
