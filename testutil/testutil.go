package testutil

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

const FooDTM com.DTM = 1

type Foo struct {
	Num int
	Str string
}

var FooSer = fooSer{}

type fooSer struct{}

func (s fooSer) Marshal(foo Foo, bs []byte) (n int) {
	n = varint.Int.Marshal(foo.Num, bs)
	n += ord.String.Marshal(foo.Str, bs[n:])
	return
}

func (s fooSer) Unmarshal(bs []byte) (foo Foo, n int, err error) {
	foo.Num, n, err = varint.Int.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	foo.Str, n1, err = ord.String.Unmarshal(bs[n:])
	n += n1
	return
}

func (s fooSer) Size(foo Foo) (size int) {
	size = varint.Int.Size(foo.Num)
	return size + ord.String.Size(foo.Str)
}

func (s fooSer) Skip(bs []byte) (n int, err error) {
	n, err = varint.Int.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.String.Skip(bs[n:])
	n += n1
	return
}
