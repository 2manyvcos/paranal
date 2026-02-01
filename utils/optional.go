package utils

import "encoding/json"

type Optional[ValueType any] struct {
	IsDefined bool
	Value     ValueType
}

func (s *Optional[ValueType]) UnmarshalJSON(d []byte) error {
	s.IsDefined = true
	return json.Unmarshal(d, &s.Value)
}

func (s *Optional[ValueType]) ApplyIfDefined(other *ValueType) {
	if s.IsDefined {
		*other = s.Value
	}
}
