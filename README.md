demo_cbor: greenpack vs cbor
=========

I took a quick look at a new kid on the block, CBOR ( https://github.com/fxamacker/cbor ). 

It was not around when I created greenpack ( https://github.com/glycerine/greenpack ). 

TL;DR: It still has a ways to go to be interesing/offer competitive performance. 

Greenpack is 4x or 400% faster on Unmarshalling, and 3x faster on Marshal.

Code generation for the win.

Importantly, for diagnostics, CBOR cannot be converted to JSON without
knowing the expected structure (what is encoded in the 
potential arbitrary data) in advance. That's a huge downside.

~~~
goos: darwin
goarch: amd64
pkg: github.com/glycerine/demo_cbor
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
Benchmark_CBOR_UnmarshalMyStruct
Benchmark_CBOR_UnmarshalMyStruct-8         	  924885	      1302 ns/op	  78.36 MB/s	      88 B/op	       7 allocs/op
Benchmark_CBOR_MarshalMsgMyStruct
Benchmark_CBOR_MarshalMsgMyStruct-8        	 1467642	       827.2 ns/op	     168 B/op	       3 allocs/op
Benchmark_Greenpack_UnmarshalMyStruct
Benchmark_Greenpack_UnmarshalMyStruct-8    	 3651930	       322.4 ns/op	 409.42 MB/s	      16 B/op	       2 allocs/op
Benchmark_Greenpack_MarshalMsgMyStruct
Benchmark_Greenpack_MarshalMsgMyStruct-8   	 4343541	       277.5 ns/op	     176 B/op	       1 allocs/op
~~~
