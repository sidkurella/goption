package option

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

// Scan implements database/sql.Scanner for Option.
// A nil source maps to Nothing. Non-nil values map to Some.
func (o *Option[T]) Scan(src any) error {
	if o == nil {
		return fmt.Errorf("option: Scan on nil pointer")
	}

	if src == nil {
		*o = Nothing[T]()
		return nil
	}

	var target T
	if scanner, ok := any(&target).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		*o = Some(target)
		return nil
	}

	sourceValue := reflect.ValueOf(src)
	for sourceValue.Kind() == reflect.Pointer {
		if sourceValue.IsNil() {
			*o = Nothing[T]()
			return nil
		}
		sourceValue = sourceValue.Elem()
	}

	targetType := reflect.TypeOf((*T)(nil)).Elem()
	if sourceValue.Type().AssignableTo(targetType) {
		*o = Some(sourceValue.Interface().(T))
		return nil
	}

	if sourceValue.Type().ConvertibleTo(targetType) {
		*o = Some(sourceValue.Convert(targetType).Interface().(T))
		return nil
	}

	return fmt.Errorf("option: cannot scan type %T into Option[%s]", src, targetType.String())
}

// Value implements database/sql/driver.Valuer for Option.
// Nothing maps to SQL NULL.
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNothing() {
		return nil, nil
	}

	v := o.Unwrap()
	if valuer, ok := any(v).(driver.Valuer); ok {
		return valuer.Value()
	}
	if valuer, ok := any(&v).(driver.Valuer); ok {
		return valuer.Value()
	}

	converted, err := driver.DefaultParameterConverter.ConvertValue(v)
	if err != nil {
		return nil, fmt.Errorf("option: cannot convert %T to driver.Value: %w", v, err)
	}
	return converted, nil
}
