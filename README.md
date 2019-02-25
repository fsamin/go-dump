# Go-Dump

Go-Dump is a package which helps you to dump a struct to `SdtOut`, any `io.Writer`, or a `map[string]string`.

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/fsamin/go-dump) [![Build Status](https://travis-ci.org/fsamin/go-dump.svg?branch=master)](https://travis-ci.org/fsamin/go-dump) [![Go Report Card](https://goreportcard.com/badge/github.com/fsamin/go-dump)](https://goreportcard.com/report/github.com/fsamin/go-dump)

## Sample usage

````golang
type T struct {
    A int
    B string
}

a := T{23, "foo bar"}

dump.FDump(out, a)
````

Will print

````bash
T.A: 23
T.B: foo bar
````

## Usage with a map

```golang
type T struct {
    A int
    B string
}

a := T{23, "foo bar"}

m, _ := dump.ToMap(a)
```

Will return such a map:

| KEY           | Value         |
| ------------- | ------------- |
| T.A           | 23            |
| T.B           | foo bar       |

## Formatting keys

```golang
    dump.ToMap(a, dump.WithDefaultLowerCaseFormatter())
```

## More examples

See [unit tests](dump_test.go) for more examples.

## Dependencies

Go-Dump needs Go >= 1.7

No external dependencies :)
