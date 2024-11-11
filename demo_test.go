package demo_cbor

import (
	"encoding/json"
	"fmt"

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
	}.EncMode()
	panicOn(err)

	dec, err := cbor.DecOptions{
		//Time: cbor.TimeRFC3339Nano, // does not seem to be needed, nor even available.
		//		ByteSliceLaterFormat: cbor.ByteSliceLaterFormatBase64,
		//		String:               cbor.StringToByteString,
		//		ByteArray:            cbor.ByteArrayToArray,
	}.DecMode()
	panicOn(err)

	v := &MyStruct{
		A: make(map[string]int),
		B: make(map[int]string),
		//C:  []float64{math.NaN(), 0, 1, -1, 43.5},
		C:  []float64{-3e-100, 0, 1, -1, 43.5},
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

	// output to JSON without knowing the content in advance? no, cbor cannot do it(?)

	// Claude suggests decoding into an empty interface as a way to do this.

	// Example usage
	//cborData := []byte{0xa1, 0x65, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x65, 0x77, 0x6f, 0x72, 0x6c, 0x64} // {"hello":"world"}

	cborData := b // JSON chokes on NaN, of course.

	jsonStr, err := CBORToJSON(cborData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(jsonStr)

}

// func CBORToJSON(cborData []byte) (string, error) {
// 	// Decode CBOR into interface{} to handle unknown structures
// 	var v interface{}
// 	if err := cbor.Unmarshal(cborData, &v); err != nil {
// 		return "", fmt.Errorf("CBOR decode error: %w", err)
// 	}

// 	// Convert to JSON
// 	jsonBytes, err := json.Marshal(v)
// 	if err != nil {
// 		return "", fmt.Errorf("JSON encode error: %w", err)
// 	}

// 	return string(jsonBytes), nil
// }
//
// barfed with: Error: JSON encode error: json: unsupported type: map[interface {}]interface {}

func convertToJSONCompatible(v interface{}) interface{} {
	switch v := v.(type) {
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range v {
			newMap[fmt.Sprint(k)] = convertToJSONCompatible(v)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, v := range v {
			newSlice[i] = convertToJSONCompatible(v)
		}
		return newSlice
	default:
		return v
	}
}

func CBORToJSON(cborData []byte) (string, error) {
	var v interface{}
	if err := cbor.Unmarshal(cborData, &v); err != nil {
		return "", fmt.Errorf("CBOR decode error: %w", err)
	}

	// Convert to JSON-compatible types
	converted := convertToJSONCompatible(v)

	jsonBytes, err := json.Marshal(converted)
	if err != nil {
		return "", fmt.Errorf("JSON encode error: %w", err)
	}

	return string(jsonBytes), nil
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
