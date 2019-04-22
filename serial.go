package serial

import (
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Encoder struct {
	w      io.Writer
	wg     sync.WaitGroup
	groups []string
}

func (e *Encoder) Encode(entity interface{}) error {
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	if err := e.encodeField(t, v, e.w); err != nil {
		return err
	}
	return nil
}
func (e *Encoder) encodeAny(entity interface{}, w io.Writer) error {
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	first := true
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if gr, ok := f.Tag.Lookup("group"); ok {
			groups := strings.Split(gr, ",")
			name := f.Name
			if tname, ok := f.Tag.Lookup("json"); ok {
				name = tname
			}
			if name != "-" && e.intersect(groups, e.groups) {
				if !first {
					w.Write([]byte(`,`))
				}
				w.Write([]byte(`"` + name + `":`))
				e.encodeField(f.Type, fv, w)
				first = false
			}
		}
	}
	return nil
}

func (e *Encoder) intersect(a []string, b []string) bool {
	for _, sa := range a {
		for _, sb := range b {
			if sa == sb {
				return true
			}
		}
	}
	return false
}

func (e *Encoder) encodeField(ft reflect.Type, v reflect.Value, w io.Writer) error {
	switch ft.Kind() {
	case reflect.Struct:
		w.Write([]byte("{"))
		if err := e.encodeAny(v.Interface(), w); err != nil {
			return err
		}
		w.Write([]byte("}"))
	case reflect.String:
		w.Write([]byte(`"` + v.String() + `"`))
	case reflect.Bool:
		if v.Bool() {
			w.Write([]byte("true"))
		} else {
			w.Write([]byte("false"))
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		w.Write([]byte(strconv.Itoa(int(v.Int()))))
	case reflect.Float32, reflect.Float64:
		w.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))

	}
	return nil
}

func (e *Encoder) AddGroup(group string) *Encoder {
	e.groups = append(e.groups, group)
	return e
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, groups: []string{}}
}
