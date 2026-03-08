package option

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON implements json.Marshaler for Option.
// Some values are marshaled as their underlying value, while Nothing marshals as null.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNothing() {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements json.Unmarshaler for Option.
// null unmarshals to Nothing. Non-null values unmarshal to Some.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return fmt.Errorf("option: UnmarshalJSON on nil pointer")
	}

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*o = Nothing[T]()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*o = Some(value)
	return nil
}
