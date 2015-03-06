package logstash

import (
	"encoding/json"
	"fmt"
)

func (attr Attribute) MarshalJSON() ([]byte, error) {
	if attr.Ival != nil {
		return []byte(fmt.Sprintf("%d", *attr.Ival)), nil
	}
	if attr.Sval != nil {
		v, err := json.Marshal(attr.Sval)
		if err != nil {
			return nil, err
		}

		return []byte(fmt.Sprintf("%s", v)), nil
	}
	if attr.Bval != nil {
		v, err := json.Marshal(attr.Bval)
		if err != nil {
			return nil, err
		}

		// The value for _bytes would naturally be true, but we want to keep
		// the literal as a pure map[string]string so we use an empty string.
		// It's just the presence of _bytes that matters anyway
		return []byte(fmt.Sprintf("{ \"value\": %s, \"_bytes\": \"\" }", string(v))), nil
	}
	if attr.Tval != nil {
		v, err := json.Marshal(attr.Tval)
		if err != nil {
			return nil, err
		}

		return []byte(fmt.Sprintf("%s", v)), nil
	}
	return []byte(fmt.Sprintf("1")), nil
}
