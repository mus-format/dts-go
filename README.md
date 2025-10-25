# dts-go

[![Go Reference](https://pkg.go.dev/badge/github.com/mus-format/dts-go.svg)](https://pkg.go.dev/github.com/mus-format/dts-go)
[![GoReportCard](https://goreportcard.com/badge/mus-format/dts-go)](https://goreportcard.com/report/github.com/mus-format/dts-go)
[![codecov](https://codecov.io/gh/mus-format/dts-go/graph/badge.svg?token=VB6E8M2PFE)](https://codecov.io/gh/mus-format/dts-go)

**dts-go** provides [DTM](https://medium.com/p/21d7be309e8d) support for the `mus-go`
serializer (DTS stands for Data Type metadata Support).

`dts-go` is particularly useful when deserializing data with an unpredictable
type. This could include completely different types, such as `Foo` and `Bar`,
or different versions of the same data, such as `FooV1` and `FooV2`.

It encodes a DTM (which is simply a number) along with the data itself,
allowing one type to be easily distinguished from another. Let’s see how:

```go
package main

import (
  "math/rand"

  com "github.com/mus-format/common-go"
  dts "github.com/mus-format/dts-go"
  "github.com/mus-format/mus-go"
)
  
type Foo struct{...}
type Bar struct{..}

// DTM (Data Type Metadata) definitions.
const (
  FooDTM com.DTM = iota + 1
  BarDTM
)

// Serializers.
var (
  FooMUS = ...
  BarMUS = ...
)

// DTS (Data Type metadata Support) definitions.
var (
  FooDTS = dts.New[Foo](FooDTM, FooMUS)
  BarDTS = dts.New[Bar](BarDTM, BarMUS)
)

func main() {
  // Make a random data and Unmarshal DTM.
  bs := randomData()
  dtm, n, err := dts.DTMSer.Unmarshal(bs)
  if err != nil {
    panic(err)
  }

  // Deserialize and process data depending on the DTM.
  switch dtm {
  case FooDTM:
    foo, _, err := FooDTS.UnmarshalData(bs[n:])
    if err != nil {
      panic(err)
    }
    // process foo ...
    fmt.Println(foo)
  case BarDTM:
    bar, _, err := BarDTS.UnmarshalData(bs[n:])
    if err != nil {
      panic(err)
    }
    // process bar ...
    fmt.Println(bar)
  default:
    panic(fmt.Sprintf("unexpected %v DTM", dtm))
  }
}
}

func randomData() (bs []byte) {...}
```

A full example can be found at [examples-go](https://github.com/mus-format/examples-go/tree/main/dts)
