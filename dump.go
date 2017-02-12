package dump

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
)

func Dump(i interface{}) error {
	return FDump(os.Stdout, i)
}

func FDump(w io.Writer, i interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	return fDumpStruct(w, i)
}

func fDumpStruct(w io.Writer, i interface{}, roots ...string) error {
	var s reflect.Value
	if reflect.ValueOf(i).Kind() == reflect.Ptr {
		s = reflect.ValueOf(i).Elem()
	} else {
		s = reflect.ValueOf(i)
	}

	switch s.Kind() {
	case reflect.Struct:
		roots = append(roots, s.Type().Name())
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			if f.Kind() == reflect.Ptr {
				f = f.Elem()
			}
			switch f.Kind() {
			case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
				if err := fDumpStruct(w, f.Interface(), append(roots, s.Type().Field(i).Name)...); err != nil {
					return err
				}
			default:
				res := fmt.Sprintf("%s.%s: %v\n", strings.Join(roots, "."), s.Type().Field(i).Name, f.Interface())
				w.Write([]byte(res))
			}
		}
	case reflect.Array, reflect.Slice:
		if err := fDumpArray(w, i, roots...); err != nil {
			return err
		}
		return nil
	case reflect.Map:
		if err := fDumpMap(w, i, roots...); err != nil {
			return err
		}
		return nil
	default:
		roots = append(roots, s.Type().Name())
		res := fmt.Sprintf("%s.%s: %v\n", strings.Join(roots, "."), s.Type().Name(), s.Interface())
		_, err := w.Write([]byte(res))
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func fDumpArray(w io.Writer, i interface{}, roots ...string) error {
	v := reflect.ValueOf(i)
	for i := 0; i < v.Len(); i++ {
		var l string
		var croots []string
		if len(roots) > 0 {
			l = roots[len(roots)-1:][0]
			croots = roots[:len(roots)-1]
		}
		croots = append(roots, fmt.Sprintf("%s%d", l, i))
		f := v.Index(i)
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}
		switch f.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
			if err := fDumpStruct(w, f.Interface(), croots...); err != nil {
				return err
			}
		default:
			res := fmt.Sprintf("%s: %v\n", strings.Join(croots, "."), f.Interface())
			w.Write([]byte(res))
		}
	}

	return nil
}

func fDumpMap(w io.Writer, i interface{}, roots ...string) error {
	v := reflect.ValueOf(i)
	keys := v.MapKeys()
	for _, k := range keys {
		roots := append(roots, fmt.Sprintf("%v", k.Interface()))
		if err := fDumpStruct(w, v.MapIndex(k).Interface(), roots...); err != nil {
			return err
		}
	}
	return nil
}

type mapWriter struct {
	data map[string]string
}

func (m *mapWriter) Write(p []byte) (int, error) {
	if m.data == nil {
		m.data = map[string]string{}
	}
	tuple := strings.Split(string(p), ":")
	if len(tuple) != 2 {
		return 0, errors.New("malformatted bytes")
	}
	tuple[1] = strings.Replace(tuple[1], "\n", "", -1)
	m.data[strings.TrimSpace(tuple[0])] = strings.TrimSpace(tuple[1])
	return len(p), nil
}

func ToMap(i interface{}) (map[string]string, error) {
	m := mapWriter{}
	err := FDump(&m, i)
	return m.data, err
}
