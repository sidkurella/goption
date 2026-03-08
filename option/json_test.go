package option_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sidkurella/goption/option"
)

type jsonTestPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestOptionJSONInterfaces(t *testing.T) {
	var _ json.Marshaler = option.Option[int]{}
	var _ json.Unmarshaler = (*option.Option[int])(nil)
}

func TestOption_MarshalJSON(t *testing.T) {
	t.Run("Some primitive", func(t *testing.T) {
		out, err := json.Marshal(option.Some(42))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(out) != "42" {
			t.Fatalf("got %s, expected 42", string(out))
		}
	})

	t.Run("Some struct", func(t *testing.T) {
		out, err := json.Marshal(option.Some(jsonTestPayload{Name: "a", Age: 3}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(out) != `{"name":"a","age":3}` {
			t.Fatalf("got %s, expected object json", string(out))
		}
	})

	t.Run("Nothing", func(t *testing.T) {
		out, err := json.Marshal(option.Nothing[int]())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(out) != "null" {
			t.Fatalf("got %s, expected null", string(out))
		}
	})
}

func TestOption_UnmarshalJSON(t *testing.T) {
	t.Run("null to Nothing", func(t *testing.T) {
		var out option.Option[int]
		if err := json.Unmarshal([]byte("null"), &out); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !out.IsNothing() {
			t.Fatalf("expected Nothing, got %v", out)
		}
	})

	t.Run("primitive to Some", func(t *testing.T) {
		var out option.Option[int]
		if err := json.Unmarshal([]byte("42"), &out); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := option.Some(42)
		if out != expected {
			t.Fatalf("got %v, expected %v", out, expected)
		}
	})

	t.Run("struct to Some", func(t *testing.T) {
		var out option.Option[jsonTestPayload]
		if err := json.Unmarshal([]byte(`{"name":"a","age":3}`), &out); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := option.Some(jsonTestPayload{Name: "a", Age: 3})
		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("got %v, expected %v", out, expected)
		}
	})

	t.Run("invalid payload returns error", func(t *testing.T) {
		var out option.Option[int]
		if err := json.Unmarshal([]byte(`"nope"`), &out); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("nil receiver returns error", func(t *testing.T) {
		var out *option.Option[int]
		if err := out.UnmarshalJSON([]byte("1")); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestOption_JSONRoundTrip(t *testing.T) {
	t.Run("Some round-trip", func(t *testing.T) {
		in := option.Some(jsonTestPayload{Name: "abc", Age: 12})
		bytes, err := json.Marshal(in)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}
		var out option.Option[jsonTestPayload]
		if err := json.Unmarshal(bytes, &out); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if !reflect.DeepEqual(out, in) {
			t.Fatalf("got %v, expected %v", out, in)
		}
	})

	t.Run("Nothing round-trip", func(t *testing.T) {
		in := option.Nothing[jsonTestPayload]()
		bytes, err := json.Marshal(in)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}
		if string(bytes) != "null" {
			t.Fatalf("got %s, expected null", string(bytes))
		}
		var out option.Option[jsonTestPayload]
		if err := json.Unmarshal(bytes, &out); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if !reflect.DeepEqual(out, in) {
			t.Fatalf("got %v, expected %v", out, in)
		}
	})
}

type jsonContainer struct {
	Present option.Option[int] `json:"present"`
	Absent  option.Option[int] `json:"absent"`
}

func TestOption_JSONInStruct(t *testing.T) {
	out, err := json.Marshal(jsonContainer{
		Present: option.Some(1),
		Absent:  option.Nothing[int](),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != `{"present":1,"absent":null}` {
		t.Fatalf("got %s, expected object with null", string(out))
	}

	var c jsonContainer
	if err := json.Unmarshal([]byte(`{"present":2,"absent":null}`), &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Present != option.Some(2) {
		t.Fatalf("present got %v, expected Some(2)", c.Present)
	}
	if !c.Absent.IsNothing() {
		t.Fatalf("absent got %v, expected Nothing", c.Absent)
	}
}

func TestOption_JSONMissingKeySemantics(t *testing.T) {
	t.Run("missing key on fresh struct leaves zero value", func(t *testing.T) {
		var c jsonContainer
		if err := json.Unmarshal([]byte(`{"present":2}`), &c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Present != option.Some(2) {
			t.Fatalf("present got %v, expected Some(2)", c.Present)
		}
		if !c.Absent.IsNothing() {
			t.Fatalf("absent got %v, expected Nothing", c.Absent)
		}
	})

	t.Run("missing key on reused struct preserves previous value", func(t *testing.T) {
		c := jsonContainer{
			Present: option.Some(1),
			Absent:  option.Some(99),
		}
		if err := json.Unmarshal([]byte(`{"present":3}`), &c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Present != option.Some(3) {
			t.Fatalf("present got %v, expected Some(3)", c.Present)
		}
		if c.Absent != option.Some(99) {
			t.Fatalf("absent got %v, expected Some(99)", c.Absent)
		}
	})
}
