package option_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

	"github.com/sidkurella/goption/option"
)

type scannedInt int

func (s *scannedInt) Scan(src any) error {
	v, ok := src.(int64)
	if !ok {
		return errors.New("expected int64")
	}
	*s = scannedInt(v + 1)
	return nil
}

type valuedString string

func (v valuedString) Value() (driver.Value, error) {
	return "prefix:" + string(v), nil
}

type valuedStringPtr string

func (v *valuedStringPtr) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	return "ptr:" + string(*v), nil
}

func TestOptionSQLInterfaces(t *testing.T) {
	var _ sql.Scanner = (*option.Option[int])(nil)
	var _ driver.Valuer = option.Option[int]{}
}

func TestOption_Scan(t *testing.T) {
	t.Run("nil source maps to Nothing", func(t *testing.T) {
		var o option.Option[int]
		if err := o.Scan(nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !o.IsNothing() {
			t.Fatalf("expected Nothing, got %v", o)
		}
	})

	t.Run("assignable source", func(t *testing.T) {
		var o option.Option[int64]
		if err := o.Scan(int64(42)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := option.Some(int64(42))
		if o != expected {
			t.Fatalf("got %v, expected %v", o, expected)
		}
	})

	t.Run("convertible source", func(t *testing.T) {
		var o option.Option[int]
		if err := o.Scan(int64(42)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := option.Some(42)
		if o != expected {
			t.Fatalf("got %v, expected %v", o, expected)
		}
	})

	t.Run("delegates to sql.Scanner", func(t *testing.T) {
		var o option.Option[scannedInt]
		if err := o.Scan(int64(5)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := option.Some(scannedInt(6))
		if o != expected {
			t.Fatalf("got %v, expected %v", o, expected)
		}
	})

	t.Run("unsupported source type", func(t *testing.T) {
		var o option.Option[int]
		if err := o.Scan(struct{}{}); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("nil receiver returns error", func(t *testing.T) {
		var o *option.Option[int]
		if err := o.Scan(int64(1)); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("typed nil pointer source maps to Nothing", func(t *testing.T) {
		var p *int64
		var o option.Option[int64]
		if err := o.Scan(p); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !o.IsNothing() {
			t.Fatalf("expected Nothing, got %v", o)
		}
	})
}

func TestOption_Value(t *testing.T) {
	t.Run("Nothing maps to nil", func(t *testing.T) {
		v, err := option.Nothing[int]().Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != nil {
			t.Fatalf("expected nil, got %v", v)
		}
	})

	t.Run("Some with primitive value", func(t *testing.T) {
		v, err := option.Some("abc").Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != "abc" {
			t.Fatalf("expected abc, got %v", v)
		}
	})

	t.Run("delegates to driver.Valuer value receiver", func(t *testing.T) {
		v, err := option.Some(valuedString("x")).Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != "prefix:x" {
			t.Fatalf("expected prefix:x, got %v", v)
		}
	})

	t.Run("delegates to driver.Valuer pointer receiver", func(t *testing.T) {
		v, err := option.Some(valuedStringPtr("y")).Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != "ptr:y" {
			t.Fatalf("expected ptr:y, got %v", v)
		}
	})

	t.Run("unsupported value type returns error", func(t *testing.T) {
		v, err := option.Some([]int{1, 2, 3}).Value()
		if err == nil {
			t.Fatalf("expected error, got value %v", v)
		}
	})

	t.Run("byte slice is accepted", func(t *testing.T) {
		b := []byte("abc")
		v, err := option.Some(b).Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out, ok := v.([]byte)
		if !ok {
			t.Fatalf("expected []byte, got %T", v)
		}
		if !reflect.DeepEqual(out, b) {
			t.Fatalf("got %v, expected %v", out, b)
		}
	})
}
