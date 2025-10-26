package dts

import (
	"testing"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/dts-go/testdata"
	"github.com/mus-format/mus-go"
	asserterror "github.com/ymz-ncnk/assert/error"
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
			asserterror.Equal(n, len(bs), t)

			afoo, n, err := fooDTS.Unmarshal(bs)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, len(bs), t)
			asserterror.EqualDeep(afoo, foo, t)

			n1, err := fooDTS.Skip(bs)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, n1, t)
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
			asserterror.Equal(n, len(bs), t)

			dtm, n, err := DTMSer.Unmarshal(bs)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(dtm, testdata.FooDTM, t)
			asserterror.Equal(n, wantDTSize, t)

			afoo, n1, err := fooDTS.UnmarshalData(bs[n:])
			asserterror.EqualError(err, nil, t)
			asserterror.EqualDeep(foo, afoo, t)
			asserterror.Equal(n1, len(bs)-wantDTSize, t)

			fooDTS.Marshal(foo, bs)
			n, err = DTMSer.Skip(bs)
			asserterror.EqualError(err, nil, t)

			n1, err = fooDTS.SkipData(bs[n:])
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n1, len(bs)-wantDTSize, t)
		})

	t.Run("DTM method should return correct DTM", func(t *testing.T) {
		var (
			fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			dtm    = fooDTS.DTM()
		)
		asserterror.Equal(dtm, testdata.FooDTM, t)
	})

	t.Run("Unamrshal should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				actualDTM = testdata.FooDTM + 3

				wantDTSize = 1
				wantErr    = com.NewWrongDTMError(testdata.FooDTM, actualDTM)

				bs     = []byte{byte(actualDTM)}
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			foo, n, err := fooDTS.Unmarshal(bs)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(foo, testdata.Foo{}, t)
			asserterror.Equal(n, wantDTSize, t)
		})

	t.Run("Skip should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				actualDTM = testdata.FooDTM + 3

				wantDTSize = 1
				wantErr    = com.NewWrongDTMError(testdata.FooDTM, actualDTM)

				dtm    = testdata.FooDTM + 3
				bs     = []byte{byte(dtm)}
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			n, err := fooDTS.Skip(bs)
			asserterror.EqualError(err, wantErr, t)
			asserterror.Equal(n, wantDTSize, t)
		})

	t.Run("If UnmarshalDTM fails with an error, Unmarshal should return it",
		func(t *testing.T) {
			var (
				wantFoo = testdata.Foo{}
				wantN   = 0
				wantErr = mus.ErrTooSmallByteSlice

				bs     = []byte{}
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			foo, n, err := fooDTS.Unmarshal(bs)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(foo, wantFoo, t)
			asserterror.Equal(n, wantN, t)
		})

	t.Run("If UnmarshalDTM fails with an error, Skip should return it",
		func(t *testing.T) {
			var (
				wantN   = 0
				wantErr = mus.ErrTooSmallByteSlice

				bs     = []byte{}
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			n, err := fooDTS.Skip(bs)
			asserterror.EqualError(err, wantErr, t)
			asserterror.Equal(n, wantN, t)
		})

	t.Run("If varint.PositiveInt.Unmarshal fails with an error, UnmarshalDTM should return it",
		func(t *testing.T) {
			var (
				wantDTM com.DTM = 0
				wantN           = 0
				wantErr         = mus.ErrTooSmallByteSlice

				bs = []byte{}
			)

			dtm, n, err := DTMSer.Unmarshal(bs)
			asserterror.EqualError(err, wantErr, t)
			asserterror.Equal(dtm, wantDTM, t)
			asserterror.Equal(n, wantN, t)
		})
}
