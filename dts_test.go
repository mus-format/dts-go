package dts

import (
	"reflect"
	"testing"

	"github.com/mus-format/dts-go/testdata"
	"github.com/mus-format/mus-go"
)

func TestDTS(t *testing.T) {
	t.Run("Marshal, Unmarshal, Size, Skip methods should work correctly",
		func(t *testing.T) {
			var (
				foo    = testdata.Foo{Num: 11, Str: "hello world"}
				fooDTS = New[testdata.Foo](testdata.FooDTM, testdata.FooSer)
				bs     = make([]byte, fooDTS.Size(foo))
			)
			n := fooDTS.Marshal(foo, bs)
			if n != len(bs) {
				t.Fatalf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			afoo, n, err := fooDTS.Unmarshal(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != len(bs) {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
			n1, err := fooDTS.Skip(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n1 != n {
				t.Errorf("unexpected n1, want '%v' actual '%v'", n, n1)
			}
		})

	t.Run("Marshal, UnmarshalDTM, UnmarshalData, Size, SkipDTM, SkipData methods should work correctly",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				foo        = testdata.Foo{Num: 11, Str: "hello world"}
				fooDTS     = New[testdata.Foo](testdata.FooDTM, testdata.FooSer)
				bs         = make([]byte, fooDTS.Size(foo))
			)
			n := fooDTS.Marshal(foo, bs)
			if n != len(bs) {
				t.Fatalf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			dtm, n, err := DTMSer.Unmarshal(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", 1, n)
			}
			if dtm != testdata.FooDTM {
				t.Errorf("unexpected dtm, want '%v' actual '%v'", testdata.FooDTM, dtm)
			}
			afoo, n1, err := fooDTS.UnmarshalData(bs[n:])
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n1 != len(bs)-wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n1)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
			fooDTS.Marshal(foo, bs)
			n, err = DTMSer.Skip(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			n1, err = fooDTS.SkipData(bs[n:])
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n1 != len(bs)-wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n1)
			}
		})

	t.Run("DTM method should return correct DTM", func(t *testing.T) {
		var (
			fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			dtm    = fooDTS.DTM()
		)
		if dtm != testdata.FooDTM {
			t.Errorf("unexpected dtm, want '%v' actual '%v'", testdata.FooDTM, dtm)
		}
	})

	t.Run("Unamrshal should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				wantDTSize  = 1
				bs          = []byte{byte(testdata.FooDTM) + 3}
				fooDTS      = New[testdata.Foo](testdata.FooDTM, nil)
				foo, n, err = fooDTS.Unmarshal(bs)
			)
			if err != ErrWrongDTM {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrWrongDTM, err)
			}
			if !reflect.DeepEqual(foo, testdata.Foo{}) {
				t.Errorf("unexpected foo, want '%v' actual '%v'", nil, foo)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", wantDTSize, n)
			}
		})

	t.Run("Skip should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				bs         = []byte{byte(testdata.FooDTM) + 3}
				fooDTS     = New[testdata.Foo](testdata.FooDTM, nil)
				n, err     = fooDTS.Skip(bs)
			)
			if err != ErrWrongDTM {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrWrongDTM, err)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", wantDTSize, n)
			}
		})

	t.Run("If UnmarshalDTM fails with an error, Unmarshal should return it",
		func(t *testing.T) {
			var (
				bs          = []byte{}
				fooDTS      = New[testdata.Foo](testdata.FooDTM, nil)
				foo, n, err = fooDTS.Unmarshal(bs)
			)
			if err != mus.ErrTooSmallByteSlice {
				t.Errorf("unexpected error, want '%v' actual '%v'",
					mus.ErrTooSmallByteSlice,
					err)
			}
			if !reflect.DeepEqual(foo, testdata.Foo{}) {
				t.Errorf("unexpected foo, want '%v' actual '%v'", nil, foo)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})

	t.Run("If UnmarshalDTM fails with an error, Skip should return it",
		func(t *testing.T) {
			var (
				bs     = []byte{}
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
				n, err = fooDTS.Skip(bs)
			)
			if err != mus.ErrTooSmallByteSlice {
				t.Errorf("unexpected error, want '%v' actual '%v'",
					mus.ErrTooSmallByteSlice,
					err)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})

	t.Run("If varint.PositiveInt.Unmarshal fails with an error, UnmarshalDTM should return it",
		func(t *testing.T) {
			var (
				bs          = []byte{}
				dtm, n, err = DTMSer.Unmarshal(bs)
			)
			if err != mus.ErrTooSmallByteSlice {
				t.Errorf("unexpected error, want '%v' actual '%v'",
					mus.ErrTooSmallByteSlice,
					err)
			}
			if dtm != 0 {
				t.Errorf("unexpected dtm, want '%v' actual '%v'", 0, dtm)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})
}
