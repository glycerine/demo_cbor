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
	A  map[string]int
	B  map[int]string
	C  []float64
	Tm time.Time
}

func TestEncode(t *testing.T) {

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
	b, err := enc.Marshal(v) // encode v to []byte b
	panicOn(err)

	var v2 MyStruct
	//err = cbor.Unmarshal(b, &v2)
	err = dec.Unmarshal(b, &v2)

	fmt.Printf("v = '%#v'\n", v)
	fmt.Printf("v2 = '%#v'\n", v2)

	// output to JSON without knowing the content in advance? no, cbor cannot do it.

}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
