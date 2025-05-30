// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fnv implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function.
//
// All the hash.Hash implementations returned by this package also
// implement encoding.BinaryMarshaler and encoding.BinaryUnmarshaler to
// marshal and unmarshal the internal state of the hash.
package fnv

import (
	"errors"
	"hash"
	"internal/byteorder"
	"math/bits"
)

type (
	sum32   uint32
	sum32a  uint32
	sum64   uint64
	sum64a  uint64
	sum128  [2]uint64
	sum128a [2]uint64
)

const (
	offset32        = 2166136261
	offset64        = 14695981039346656037
	offset128Lower  = 0x62b821756295c58d
	offset128Higher = 0x6c62272e07bb0142
	prime32         = 16777619
	prime64         = 1099511628211
	prime128Lower   = 0x13b
	prime128Shift   = 24
)

// New32 returns a new 32-bit FNV-1 [hash.Hash].
// Its Sum method will lay the value out in big-endian byte order.
func New32() hash.Hash32 {
	var s sum32 = offset32
	return &s
}

// New32a returns a new 32-bit FNV-1a [hash.Hash].
// Its Sum method will lay the value out in big-endian byte order.
func New32a() hash.Hash32 {
	var s sum32a = offset32
	return &s
}

// New64 returns a new 64-bit FNV-1 [hash.Hash].
// Its Sum method will lay the value out in big-endian byte order.
func New64() hash.Hash64 {
	var s sum64 = offset64
	return &s
}

// New64a returns a new 64-bit FNV-1a [hash.Hash].
// Its Sum method will lay the value out in big-endian byte order.
func New64a() hash.Hash64 {
	var s sum64a = offset64
	return &s
}

// New128 returns a new 128-bit FNV-1 [hash.Hash].
// Its Sum method will lay the value out in big-endian byte order.
func New128() hash.Hash {
	var s sum128
	s[0] = offset128Higher
	s[1] = offset128Lower
	return &s
}

// New128a returns a new 128-bit FNV-1a [hash.Hash].
// Its Sum method will lay the value out in big-endian byte order.
func New128a() hash.Hash {
	var s sum128a
	s[0] = offset128Higher
	s[1] = offset128Lower
	return &s
}

func (s *sum32) Reset()   { *s = offset32 }
func (s *sum32a) Reset()  { *s = offset32 }
func (s *sum64) Reset()   { *s = offset64 }
func (s *sum64a) Reset()  { *s = offset64 }
func (s *sum128) Reset()  { s[0] = offset128Higher; s[1] = offset128Lower }
func (s *sum128a) Reset() { s[0] = offset128Higher; s[1] = offset128Lower }

func (s *sum32) Sum32() uint32  { return uint32(*s) }
func (s *sum32a) Sum32() uint32 { return uint32(*s) }
func (s *sum64) Sum64() uint64  { return uint64(*s) }
func (s *sum64a) Sum64() uint64 { return uint64(*s) }

func (s *sum32) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash *= prime32
		hash ^= sum32(c)
	}
	*s = hash
	return len(data), nil
}

func (s *sum32a) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash ^= sum32a(c)
		hash *= prime32
	}
	*s = hash
	return len(data), nil
}

func (s *sum64) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash *= prime64
		hash ^= sum64(c)
	}
	*s = hash
	return len(data), nil
}

func (s *sum64a) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash ^= sum64a(c)
		hash *= prime64
	}
	*s = hash
	return len(data), nil
}

func (s *sum128) Write(data []byte) (int, error) {
	for _, c := range data {
		// Compute the multiplication
		s0, s1 := bits.Mul64(prime128Lower, s[1])
		s0 += s[1]<<prime128Shift + prime128Lower*s[0]
		// Update the values
		s[1] = s1
		s[0] = s0
		s[1] ^= uint64(c)
	}
	return len(data), nil
}

func (s *sum128a) Write(data []byte) (int, error) {
	for _, c := range data {
		s[1] ^= uint64(c)
		// Compute the multiplication
		s0, s1 := bits.Mul64(prime128Lower, s[1])
		s0 += s[1]<<prime128Shift + prime128Lower*s[0]
		// Update the values
		s[1] = s1
		s[0] = s0
	}
	return len(data), nil
}

func (s *sum32) Size() int   { return 4 }
func (s *sum32a) Size() int  { return 4 }
func (s *sum64) Size() int   { return 8 }
func (s *sum64a) Size() int  { return 8 }
func (s *sum128) Size() int  { return 16 }
func (s *sum128a) Size() int { return 16 }

func (s *sum32) BlockSize() int   { return 1 }
func (s *sum32a) BlockSize() int  { return 1 }
func (s *sum64) BlockSize() int   { return 1 }
func (s *sum64a) BlockSize() int  { return 1 }
func (s *sum128) BlockSize() int  { return 1 }
func (s *sum128a) BlockSize() int { return 1 }

func (s *sum32) Sum(in []byte) []byte {
	v := uint32(*s)
	return byteorder.BEAppendUint32(in, v)
}

func (s *sum32a) Sum(in []byte) []byte {
	v := uint32(*s)
	return byteorder.BEAppendUint32(in, v)
}

func (s *sum64) Sum(in []byte) []byte {
	v := uint64(*s)
	return byteorder.BEAppendUint64(in, v)
}

func (s *sum64a) Sum(in []byte) []byte {
	v := uint64(*s)
	return byteorder.BEAppendUint64(in, v)
}

func (s *sum128) Sum(in []byte) []byte {
	ret := byteorder.BEAppendUint64(in, s[0])
	return byteorder.BEAppendUint64(ret, s[1])
}

func (s *sum128a) Sum(in []byte) []byte {
	ret := byteorder.BEAppendUint64(in, s[0])
	return byteorder.BEAppendUint64(ret, s[1])
}

const (
	magic32          = "fnv\x01"
	magic32a         = "fnv\x02"
	magic64          = "fnv\x03"
	magic64a         = "fnv\x04"
	magic128         = "fnv\x05"
	magic128a        = "fnv\x06"
	marshaledSize32  = len(magic32) + 4
	marshaledSize64  = len(magic64) + 8
	marshaledSize128 = len(magic128) + 8*2
)

func (s *sum32) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, magic32...)
	b = byteorder.BEAppendUint32(b, uint32(*s))
	return b, nil
}

func (s *sum32) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, marshaledSize32))
}

func (s *sum32a) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, magic32a...)
	b = byteorder.BEAppendUint32(b, uint32(*s))
	return b, nil
}

func (s *sum32a) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, marshaledSize32))
}

func (s *sum64) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, magic64...)
	b = byteorder.BEAppendUint64(b, uint64(*s))
	return b, nil
}

func (s *sum64) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, marshaledSize64))
}

func (s *sum64a) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, magic64a...)
	b = byteorder.BEAppendUint64(b, uint64(*s))
	return b, nil
}

func (s *sum64a) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, marshaledSize64))
}

func (s *sum128) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, magic128...)
	b = byteorder.BEAppendUint64(b, s[0])
	b = byteorder.BEAppendUint64(b, s[1])
	return b, nil
}

func (s *sum128) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, marshaledSize128))
}

func (s *sum128a) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, magic128a...)
	b = byteorder.BEAppendUint64(b, s[0])
	b = byteorder.BEAppendUint64(b, s[1])
	return b, nil
}

func (s *sum128a) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, marshaledSize128))
}

func (s *sum32) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic32) || string(b[:len(magic32)]) != magic32 {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize32 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum32(byteorder.BEUint32(b[4:]))
	return nil
}

func (s *sum32a) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic32a) || string(b[:len(magic32a)]) != magic32a {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize32 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum32a(byteorder.BEUint32(b[4:]))
	return nil
}

func (s *sum64) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic64) || string(b[:len(magic64)]) != magic64 {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize64 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum64(byteorder.BEUint64(b[4:]))
	return nil
}

func (s *sum64a) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic64a) || string(b[:len(magic64a)]) != magic64a {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize64 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	*s = sum64a(byteorder.BEUint64(b[4:]))
	return nil
}

func (s *sum128) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic128) || string(b[:len(magic128)]) != magic128 {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize128 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	s[0] = byteorder.BEUint64(b[4:])
	s[1] = byteorder.BEUint64(b[12:])
	return nil
}

func (s *sum128a) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic128a) || string(b[:len(magic128a)]) != magic128a {
		return errors.New("hash/fnv: invalid hash state identifier")
	}
	if len(b) != marshaledSize128 {
		return errors.New("hash/fnv: invalid hash state size")
	}
	s[0] = byteorder.BEUint64(b[4:])
	s[1] = byteorder.BEUint64(b[12:])
	return nil
}

func (d *sum32) Clone() (hash.Cloner, error) {
	r := *d
	return &r, nil
}

func (d *sum32a) Clone() (hash.Cloner, error) {
	r := *d
	return &r, nil
}

func (d *sum64) Clone() (hash.Cloner, error) {
	r := *d
	return &r, nil
}

func (d *sum64a) Clone() (hash.Cloner, error) {
	r := *d
	return &r, nil
}

func (d *sum128) Clone() (hash.Cloner, error) {
	r := *d
	return &r, nil
}

func (d *sum128a) Clone() (hash.Cloner, error) {
	r := *d
	return &r, nil
}
