package demo_cbor

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	//"log"
	"math"
	"testing"
	"time"

	"github.com/fxamacker/cbor/v2"
)

type MyStruct struct {
	A  map[string]int `zid:"0"`
	B  map[int]string `zid:"1"`
	C  []float64      `zid:"2"`
	Tm time.Time      `zid:"3"`
}

//go:generate greenpack

func TestEncode(t *testing.T) {

	// We have to use EncOptions or else the sub-second precision
	// of the timestamps is lost. That would be a non-starter.
	enc, err := cbor.EncOptions{
		Time: cbor.TimeRFC3339Nano,
		//		ByteSliceLaterFormat: cbor.ByteSliceLaterFormatBase64,
		//		String:               cbor.StringToByteString,
		//		ByteArray:            cbor.ByteArrayToArray,
	}.EncMode()
	panicOn(err)

	dec, err := cbor.DecOptions{
		//Time: cbor.TimeRFC3339Nano,
		//		ByteSliceLaterFormat: cbor.ByteSliceLaterFormatBase64,
		//		String:               cbor.StringToByteString,
		//		ByteArray:            cbor.ByteArrayToArray,
	}.DecMode()
	panicOn(err)

	v := &MyStruct{
		A:  make(map[string]int),
		B:  make(map[int]string),
		C:  []float64{math.NaN(), 0, 1, -1, 43.5},
		Tm: time.Now(),
	}
	v.A["hello"] = 34
	v.B[34] = "hello back"

	//b, err := cbor.Marshal(v) // encode v to []byte b
	b, err := enc.Marshal(&v) // encode v to []byte b
	panicOn(err)

	var v2 MyStruct
	//err = cbor.Unmarshal(b, &v2)
	err = dec.Unmarshal(b, &v2)

	fmt.Printf("v = '%#v'\n", v)
	fmt.Printf("v2 = '%#v'\n", v2)

	// output to JSON without knowing the content in advance? no, cbor cannot do it.

}

func getTestMyStruct() *MyStruct {
	v := &MyStruct{
		A:  make(map[string]int),
		B:  make(map[int]string),
		C:  []float64{math.NaN(), 0, 1, -1, 43.5},
		Tm: time.Now(),
	}
	v.A["hello"] = 34
	v.B[34] = "hello back"
	return v
}

func Benchmark_CBOR_UnmarshalMyStruct(b *testing.B) {

	v := getTestMyStruct()

	// We have to use EncOptions or else the sub-second precision
	// of the timestamps is lost. That would be a non-starter.
	enc, err := cbor.EncOptions{
		Time: cbor.TimeRFC3339Nano,
		//		ByteSliceLaterFormat: cbor.ByteSliceLaterFormatBase64,
		//		String:               cbor.StringToByteString,
		//		ByteArray:            cbor.ByteArrayToArray,
	}.EncMode()
	panicOn(err)

	dec, err := cbor.DecOptions{}.DecMode()
	panicOn(err)

	bts, _ := enc.Marshal(&v) // encode v to []byte b

	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := dec.Unmarshal(bts, &v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_CBOR_MarshalMsgMyStruct(b *testing.B) {

	v := getTestMyStruct()

	// We have to use EncOptions or else the sub-second precision
	// of the timestamps is lost. That would be a non-starter.
	enc, err := cbor.EncOptions{
		Time: cbor.TimeRFC3339Nano,
	}.EncMode()
	panicOn(err)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.Marshal(&v)
	}
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func Benchmark_Greenpack_UnmarshalMyStruct(b *testing.B) {
	v := getTestMyStruct()
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Greenpack_MarshalMsgMyStruct(b *testing.B) {
	v := getTestMyStruct()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.MarshalMsg(nil)
	}
}
