package hash2xml

import (
	"time"
)

// hashEncoder encodes the items in a hash by iterating and encoding each entry
func hashEncoder(s *Serializer, raw interface{}, path string, key ...string) (bool, error) {
	switch hash := raw.(type) {
	case map[string]interface{}:
		if len(key) > 0 {
			s.WriteStartTag(key[0])
			s.Newline()
			s.Indent()
		}

		// recursively serialize the hash
		for k, v := range hash {
			err := s.Convert(v, path, k)
			if err != nil {
				return false, err
			}
		}

		if len(key) > 0 {
			s.Dedent()
			s.WriteIndentation()
			s.WriteEndTag(key[0])
		}
		return true, nil

	default:
		return false, nil
	}
}

// arrayEncoder encodes the items in an array by iterating and encoding each entry
func arrayEncoder(s *Serializer, raw interface{}, path string, key ...string) (bool, error) {
	switch array := raw.(type) {
	case []interface{}:
		if len(key) > 0 {
			s.WriteStartTag(key[0])
			s.Newline()
			s.Indent()
		}

		// iterate the array and serialize all the values
		for _, value := range array {
			err := s.Convert(value, path)
			if err != nil {
				return false, err
			}
		}

		if len(key) > 0 {
			s.Dedent()
			s.WriteIndentation()
			s.WriteEndTag(key[0])
		}
		return true, nil

	default:
		return false, nil
	}
}

// scalarEncoder encodes scalar values to string using %s formatting
func scalarEncoder(s *Serializer, raw interface{}, path string, key ...string) (bool, error) {
	switch value := raw.(type) {
	case string, float64, bool, int, int32, int64, float32:
		var defaultKey string

		if len(key) > 0 {
			s.WriteStartTag(key[0])
		} else {
			defaultKey = s.getDefaultKey(value)
			s.WriteStartTag(defaultKey)
		}

		s.WriteScalar(value)

		if len(key) > 0 {
			s.WriteEndTag(key[0])
		} else {
			s.WriteEndTag(defaultKey)
		}
		return true, nil

	default:
		return false, nil
	}
}

// timeEncoder encodes time to a decent format
func timeEncoder(s *Serializer, raw interface{}, path string, key ...string) (bool, error) {
	switch value := raw.(type) {
	case time.Time:
		var defaultKey string

		if len(key) > 0 {
			s.WriteStartTag(key[0])
		} else {
			defaultKey = s.getDefaultKey(value)
			s.WriteStartTag(defaultKey)
		}

		// encode and write the time
		s.WriteScalar(value.Format(time.RFC3339))

		if len(key) > 0 {
			s.WriteEndTag(key[0])
		} else {
			s.WriteEndTag(defaultKey)
		}
		return true, nil

	default:
		return false, nil
	}
}
