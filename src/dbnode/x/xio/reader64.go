package xio

import (
	"encoding/binary"
	"io"
)

type BytesReader64 struct {
	data  []byte
	index int
}

func NewBytesReader64(data []byte) *BytesReader64 {
	return &BytesReader64{data: data}
}

func (r *BytesReader64) Read64() (word uint64, n byte, err error) {
	if r.index+8 <= len(r.data) {
		// NB: this compiles to a single 64 bit load followed by
		// a BSWAPQ on amd64 gc 1.13 (https://godbolt.org/z/oTK1jx).
		res := binary.BigEndian.Uint64(r.data[r.index:])
		r.index += 8
		return res, 8, nil
	}
	if r.index >= len(r.data) {
		return 0, 0, io.EOF
	}
	var res uint64
	var bytes byte
	for ; r.index < len(r.data); r.index++ {
		res = (res << 8) | uint64(r.data[r.index])
		bytes++
	}
	return res << (64 - 8*bytes), bytes, nil
}

func (r *BytesReader64) Peek64() (word uint64, n byte, err error) {
	if r.index+8 <= len(r.data) {
		// NB: this compiles to a single 64 bit load followed by
		// BSWAPQ on amd64 gc 1.13 (https://godbolt.org/z/oTK1jx).
		res := binary.BigEndian.Uint64(r.data[r.index:])
		return res, 8, nil
	}

	if r.index >= len(r.data) {
		return 0, 0, io.EOF
	}

	var res uint64
	var bytes byte
	for i := r.index; i < len(r.data); i++ {
		res = (res << 8) | uint64(r.data[i])
		bytes++
	}
	return res << (64 - 8*bytes), bytes, nil
}

func (r *BytesReader64) Reset(data []byte) {
	r.data = data
	r.index = 0
}
