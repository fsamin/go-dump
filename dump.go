package dump

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/structs"
)

// KeyFormatterFunc is a type for key formatting
type KeyFormatterFunc func(s string) string

// WithLowerCaseFormatter formats keys in lowercase
func WithLowerCaseFormatter() KeyFormatterFunc {
	return func(s string) string {
		return strings.ToLower(s)
	}
}

// WithDefaultLowerCaseFormatter formats keys in lowercase and apply default formatting
func WithDefaultLowerCaseFormatter() KeyFormatterFunc {
	f := WithDefaultFormatter()
	return func(s string) string {
		return strings.ToLower(f(s))
	}
}

// WithDefaultFormatter is the default formatter
func WithDefaultFormatter() KeyFormatterFunc {
	return func(s string) string {
		s = strings.Replace(s, " ", "_", -1)
		s = strings.Replace(s, "/", "_", -1)
		s = strings.Replace(s, ":", "_", -1)
		return s
	}
}

// NoFormatter doesn't do anything, so to be sure to avoid keys formatting, use only this formatter
func NoFormatter() KeyFormatterFunc {
	return func(s string) string {
		return s
	}
}

// Dump displays the passed parameter properties to standard out such as complete types and all
// pointer addresses used to indirect to the final value.
// See Fdump if you would prefer dumping to an arbitrary io.Writer or Sdump to
// get the formatted result as a string.
func Dump(i interface{}, formatters ...KeyFormatterFunc) error {
	if formatters == nil {
		formatters = []KeyFormatterFunc{WithDefaultFormatter()}
	}
	return Fdump(os.Stdout, i, formatters...)
}

// Fdump formats and displays the passed arguments to io.Writer w. It formats exactly the same as Dump.
func Fdump(w io.Writer, i interface{}, formatters ...KeyFormatterFunc) (err error) {
	if formatters == nil {
		formatters = []KeyFormatterFunc{WithDefaultFormatter()}
	}

	res, err := ToStringMap(i, formatters...)
	if err != nil {
		return
	}

	keys := []string{}
	for k := range res {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		var err error
		if res[k] == "" {
			_, err = fmt.Fprintf(w, "%s:\n", k)
		} else {
			_, err = fmt.Fprintf(w, "%s: %s\n", k, res[k])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Sdump returns a string with the passed arguments formatted exactly the same as Dump.
func Sdump(i interface{}, formatters ...KeyFormatterFunc) (string, error) {
	if formatters == nil {
		formatters = []KeyFormatterFunc{WithDefaultFormatter()}
	}
	m, err := ToStringMap(i, formatters...)
	if err != nil {
		return "", err
	}
	res := ""

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		res += fmt.Sprintf("%s: %s\n", k, m[k])
	}
	return res, nil
}

func valueFromInterface(i interface{}) reflect.Value {
	var f reflect.Value
	if reflect.ValueOf(i).Kind() == reflect.Ptr {
		f = reflect.ValueOf(i).Elem()
	} else {
		f = reflect.ValueOf(i)
	}

	if f.Kind() == reflect.Interface {
		if reflect.ValueOf(f.Interface()).Kind() == reflect.Ptr {
			f = reflect.ValueOf(f.Interface()).Elem()
		} else {
			f = reflect.ValueOf(f.Interface())
		}
	}
	return f
}

func fdumpInterface(w map[string]interface{}, i interface{}, roots []string, formatters ...KeyFormatterFunc) error {
	f := valueFromInterface(i)
	if !validAndNotEmpty(f) {
		k := fmt.Sprintf("%s", strings.Join(sliceFormat(roots, formatters), "."))
		w[k] = ""
		return nil
	}

	switch f.Kind() {
	case reflect.Struct:
		nodeType := append(roots, "__Type__")
		nodeTypeFormatted := strings.Join(sliceFormat(nodeType, formatters), ".")
		w[nodeTypeFormatted] = f.Type().Name()
		croots := roots
		if len(roots) == 0 {
			croots = append(roots, f.Type().Name())
		}
		if err := fdumpStruct(w, f, croots, formatters...); err != nil {
			return err
		}
	case reflect.Array, reflect.Slice:
		nodeType := append(roots, "__Type__")
		nodeTypeFormatted := strings.Join(sliceFormat(nodeType, formatters), ".")
		w[nodeTypeFormatted] = "Array"
		if err := fDumpArray(w, i, roots, formatters...); err != nil {
			return err
		}
		return nil
	case reflect.Map:
		nodeType := append(roots, "__Type__")
		nodeTypeFormatted := strings.Join(sliceFormat(nodeType, formatters), ".")
		w[nodeTypeFormatted] = "Map"
		if err := fDumpMap(w, i, roots, formatters...); err != nil {
			return err
		}
		return nil
	default:
		k := fmt.Sprintf("%s", strings.Join(sliceFormat(roots, formatters), "."))
		w[k] = f.Interface()
	}
	return nil
}

func fDumpArray(w map[string]interface{}, i interface{}, roots []string, formatters ...KeyFormatterFunc) error {
	v := reflect.ValueOf(i)

	nodeLen := append(roots, "__Len__")
	nodeLenFormatted := strings.Join(sliceFormat(nodeLen, formatters), ".")
	w[nodeLenFormatted] = v.Len()

	for i := 0; i < v.Len(); i++ {
		var l string
		//croots := roots
		var croots []string
		if len(roots) > 0 {
			l = roots[len(roots)-1:][0]
			croots = append(roots, fmt.Sprintf("%s%d", l, i))
		} else {
			croots = append(roots, fmt.Sprintf("%s%d", l, i))
		}
		f := v.Index(i)
		if err := fdumpInterface(w, f.Interface(), croots, formatters...); err != nil {
			return err
		}
	}

	return nil
}

func fDumpMap(w map[string]interface{}, i interface{}, roots []string, formatters ...KeyFormatterFunc) error {
	v := reflect.ValueOf(i)

	keys := v.MapKeys()

	nodeLen := append(roots, "__Len__")
	nodeLenFormatted := strings.Join(sliceFormat(nodeLen, formatters), ".")
	w[nodeLenFormatted] = len(keys)

	for _, k := range keys {
		key := fmt.Sprintf("%v", k.Interface())
		croots := append(roots, key)
		value := v.MapIndex(k)

		f := valueFromInterface(value.Interface())

		if validAndNotEmpty(f) && f.Type().Kind() == reflect.Struct {
			croots = append(croots, f.Type().Name())
		}

		if err := fdumpInterface(w, value.Interface(), croots, formatters...); err != nil {
			return err
		}
	}
	return nil
}

func fdumpStruct(w map[string]interface{}, s reflect.Value, roots []string, formatters ...KeyFormatterFunc) error {
	nodeLen := append(roots, "__Len__")
	nodeLenFormatted := strings.Join(sliceFormat(nodeLen, formatters), ".")
	w[nodeLenFormatted] = len(structs.Fields(s))

	structKey := fmt.Sprintf("%s", strings.Join(sliceFormat(roots, formatters), "."))
	if s.CanInterface() {
		w[structKey] = s.Interface()
	}

	for i := 0; i < s.NumField(); i++ {
		if !s.Field(i).CanInterface() {
			continue
		}
		croots := append(roots, s.Type().Field(i).Name)
		if err := fdumpInterface(w, s.Field(i).Interface(), croots, formatters...); err != nil {
			return err
		}
	}
	return nil
}

// ToStringMap formats the argument as a map[string]string. It formats exactly the same as Dump.
func ToStringMap(i interface{}, formatters ...KeyFormatterFunc) (res map[string]string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
		}
	}()
	ires := map[string]interface{}{}
	if err = fdumpInterface(ires, i, nil, formatters...); err != nil {
		return
	}
	res = map[string]string{}
	for k, v := range ires {
		vv := reflect.ValueOf(v)
		if vv.Kind() != reflect.Struct {
			res[k] = fmt.Sprintf("%s", fmt.Sprintf("%v", v))
		}
	}
	return
}

// ToMap dumps argument as a map[string]interface{}, it can also format remove some fields
func ToMap(i interface{}, filterFields []string, formatters ...KeyFormatterFunc) (res map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
		}
	}()
	res = map[string]interface{}{}
	if err = fdumpInterface(res, i, nil, formatters...); err != nil {
		return
	}
	for _, f := range filterFields {
		for k := range res {
			if strings.HasSuffix(strings.ToUpper(k), "."+strings.ToUpper(f)) {
				delete(res, k)
			}
		}
	}
	return
}

func validAndNotEmpty(v reflect.Value) bool {
	if v.IsValid() && v.CanInterface() {
		if v.Kind() == reflect.String {
			return v.String() != ""
		}
		return true
	}
	return false
}

func sliceFormat(s []string, formatters []KeyFormatterFunc) []string {
	for i := range s {
		s[i] = format(s[i], formatters)
	}
	return s
}

func format(s string, formatters []KeyFormatterFunc) string {
	for _, f := range formatters {
		s = f(s)
	}
	return s
}

// MustSdump is a helper that wraps a call to a function returning (string, error)
// and panics if the error is non-nil.
func MustSdump(i interface{}, formatters ...KeyFormatterFunc) string {
	s, err := Sdump(i, formatters...)
	if err != nil {
		panic(err)
	}
	return s
}
