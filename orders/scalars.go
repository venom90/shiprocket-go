package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type FlexibleString string
type FlexibleInt int64
type FlexibleFloat float64
type FlexibleBool bool

func (v FlexibleString) String() string {
	return string(v)
}

func (v FlexibleString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(v))
}

func (v *FlexibleString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*v = ""
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*v = FlexibleString(asString)
		return nil
	}

	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch typed := raw.(type) {
	case float64:
		*v = FlexibleString(strconv.FormatFloat(typed, 'f', -1, 64))
	case bool:
		*v = FlexibleString(strconv.FormatBool(typed))
	default:
		return fmt.Errorf("unsupported flexible string value %T", raw)
	}

	return nil
}

func (v FlexibleInt) Int64() int64 {
	return int64(v)
}

func (v FlexibleInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(v))
}

func (v *FlexibleInt) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte(`""`)) {
		*v = 0
		return nil
	}

	var asInt int64
	if err := json.Unmarshal(data, &asInt); err == nil {
		*v = FlexibleInt(asInt)
		return nil
	}

	var asFloat float64
	if err := json.Unmarshal(data, &asFloat); err == nil {
		*v = FlexibleInt(int64(asFloat))
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if strings.TrimSpace(asString) == "" {
			*v = 0
			return nil
		}
		parsed, err := strconv.ParseInt(asString, 10, 64)
		if err != nil {
			return err
		}
		*v = FlexibleInt(parsed)
		return nil
	}

	return fmt.Errorf("unsupported flexible int value %s", string(data))
}

func (v FlexibleFloat) Float64() float64 {
	return float64(v)
}

func (v FlexibleFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(v))
}

func (v *FlexibleFloat) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte(`""`)) {
		*v = 0
		return nil
	}

	var asFloat float64
	if err := json.Unmarshal(data, &asFloat); err == nil {
		*v = FlexibleFloat(asFloat)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if strings.TrimSpace(asString) == "" {
			*v = 0
			return nil
		}
		parsed, err := strconv.ParseFloat(asString, 64)
		if err != nil {
			return err
		}
		*v = FlexibleFloat(parsed)
		return nil
	}

	return fmt.Errorf("unsupported flexible float value %s", string(data))
}

func (v FlexibleBool) Bool() bool {
	return bool(v)
}

func (v FlexibleBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(v))
}

func (v *FlexibleBool) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte(`""`)) {
		*v = false
		return nil
	}

	var asBool bool
	if err := json.Unmarshal(data, &asBool); err == nil {
		*v = FlexibleBool(asBool)
		return nil
	}

	var asInt int64
	if err := json.Unmarshal(data, &asInt); err == nil {
		*v = FlexibleBool(asInt != 0)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		switch strings.ToLower(strings.TrimSpace(asString)) {
		case "", "0", "false", "no":
			*v = false
			return nil
		case "1", "true", "yes":
			*v = true
			return nil
		default:
			return fmt.Errorf("unsupported flexible bool string %q", asString)
		}
	}

	return fmt.Errorf("unsupported flexible bool value %s", string(data))
}
