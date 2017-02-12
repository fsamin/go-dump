# Go-Dump

Go-Dump is a package which helps you to dump a struct to `SdtOut`, any `io.Writer`, or a `map[string]string{}`.

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)]
(http://godoc.org/github.com/fsamin/go-dump)

[![Build Status](https://travis-ci.org/fsamin/go-dump.svg?branch=master)](https://travis-ci.org/fsamin/go-dump)

## Sample usage

````golang
type T struct {
    A int
    B string
}

a := T{23, "foo bar"}

dump.FDump(out, a)
````

Will prints

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

Will returns such a map:

| KEY           | Value         |
| ------------- | ------------- |
| T.A           | 23            |
| T.B           | foo bar       |


## More examples

See [unit tests](test/dump_test.go) for more examples.
