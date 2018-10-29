package util

import (
	"reflect"
	"testing"
	"unsafe"
)

var (
	bs []byte
)

// DANDER: SHOULD NOT read or write cap of the slice
func quickStringByte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// returns &s[0], which is not allowed in go
func stringPointer(s string) unsafe.Pointer {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Pointer(p.Data)
}

// returns &b[0], which is not allowed in go
func bytePointer(b []byte) unsafe.Pointer {
	p := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(p.Data)
}

func Benchmark_StringByte(b *testing.B) {
	s := "test"
	for i := 0; i < b.N; i++ {
		_ = Slice(s)
	}
}

func Benchmark_StringByte_Simple(b *testing.B) {
	s := "test"
	for i := 0; i < b.N; i++ {
		bs = []byte(s)
	}
}

func Benchmark_StringByte_Quick(b *testing.B) {
	s := "test"
	for i := 0; i < b.N; i++ {
		_ = quickStringByte(s)
	}

	// 根据 StringHeader 结构，尝试了下手动获取 Len 的值，感觉好麻烦
	// type StringHeader struct {
	//     Data uintptr
	//     Len  int
	// }
	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	b.Logf("ptr is %p %p, %p \n", &s, sHeader, &sHeader.Data)
	l := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(sHeader)) + 8))
	b.Log("len is", sHeader.Len, *l)
}

func Benchmark_ByteToString(b *testing.B) {
	bs := []byte("test")
	for i := 0; i < b.N; i++ {
		_ = String(bs)
	}
}
