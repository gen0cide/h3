package h3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Rubric represents a scoring rubric for a specific skill tree element
type Rubric map[Level]string

// EmptyRubric initializes an empty rubric type
func EmptyRubric() Rubric {
	r := Rubric{}
	for x := range ValidLevels {
		r[x] = "__DEFINE_CRITERIA__"
	}
	return r
}

// MarshalJSON implements the JSONMarshaler for the Rubric type
func (r Rubric) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("{")
	length := len(r)
	count := 0
	keys := []int{}
	for l := range r {
		keys = append(keys, int(l))
	}
	sort.Ints(keys)
	for i := range keys {
		l := Level(i)
		d := r[l]
		jsonVal, err := json.Marshal(d)
		if err != nil {
			return nil, err
		}
		jsonKey, err := json.Marshal(l)
		if err != nil {
			return nil, err
		}
		buf.WriteString((fmt.Sprintf("%s: %s", jsonKey, jsonVal)))
		count++
		if count < length {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// UnmarshalJSON implements the JSONUnmarshaler for the Rubric type
func (r Rubric) UnmarshalJSON(b []byte) error {
	var rubraw map[string]string
	err := json.Unmarshal(b, &rubraw)
	if err != nil {
		return err
	}
	for k, v := range rubraw {
		lev, err := LevelFromWord(k)
		if err != nil {
			return err
		}
		r[lev] = v
	}
	return nil
}
