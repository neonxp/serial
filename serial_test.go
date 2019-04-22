package serial

import (
	"bytes"
	"testing"
)

type t1 struct {
	FieldA string `json:"field_a" group:"group_1"`
	FieldB int    `json:"field_b" group:"group_2"`
	FieldC bool   `json:"field_c" group:"group_1"`
	FieldD t2     `json:"field_d" group:"group_2,group_1"`
}
type t2 struct {
	FieldA string  `json:"field_a" group:"group_2"`
	FieldB int     `json:"field_b" group:"group_2,group_1"`
	FieldC bool    `json:"field_c" group:"group_2"`
	FieldD float64 `json:"field_d" group:"group_1"`
}

func TestSerialize(t *testing.T) {
	o1 := t1{
		FieldA: "foo",
		FieldB: 123,
		FieldC: true,
		FieldD: t2{FieldA: "bar", FieldB: 321, FieldC: false, FieldD: 123.2121},
	}
	result := bytes.NewBuffer([]byte{})
	if err := NewEncoder(result).AddGroup("group_1").Encode(o1); err != nil {
		t.Error(err)
	}
	expected := `{"field_a":"foo","field_c":true,"field_d":{"field_b":321,"field_d":123.2121}}`
	if result.String() != expected {
		t.Errorf("Expected: %s got %s", expected, result.String())
	}
}

func TestDeserialize(t *testing.T) {
}
