package main

import "encoding/binary"

func parse(X BitString, r int) ([][]uint64, []byte) {
	if r <= 0 || r%64 != 0 {
		panic("parse: r must be a positive multiple of 64")
	}
	wordsPerBlock := r / 64
	l := X.Bits / r
	remBits := X.Bits % r

	out := make([][]uint64, 0, l)
	for i := 0; i < l; i++ {
		start := i * wordsPerBlock
		end := start + wordsPerBlock
		out = append(out, X.Words[start:end])
	}
	if remBits == 0 {
		return out, nil
	}

	startWord := l * wordsPerBlock
	remWords := (remBits + 63) / 64
	buf := make([]byte, remWords*8)
	for i := 0; i < remWords; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], X.Words[startWord+i])
	}
	nBytes := (remBits + 7) / 8
	return out, buf[:nBytes]
}

func pad(X []byte, rBits int) []uint64 {

	rBytes := rBits / 8

	block := make([]byte, rBytes)
	copy(block, X)

	if len(X) < rBytes {
		block[len(X)] = 0x80
	}

	words := make([]uint64, rBytes/8)
	for i := 0; i < len(words); i++ {
		words[i] = binary.BigEndian.Uint64(block[i*8 : (i+1)*8])
	}
	return words

}

func fromByteToUint64(b []byte, l int) []uint64 {
	if l < 0 {
		panic("l must be >= 0")
	}
	if l == 0 {
		return []uint64{}
	}

	totalBits := len(b) * 8
	if l > totalBits {
		panic("l exceeds available bits in b")
	}

	needBytes := (l + 7) / 8

	tmp := make([]byte, needBytes)
	copy(tmp, b[:needBytes])

	rem := l % 8
	if rem != 0 {
		mask := byte(0xFF) << (8 - rem)
		tmp[needBytes-1] &= mask
	}

	nWords := (needBytes + 7) / 8
	out := make([]uint64, nWords)

	for i := 0; i < nWords; i++ {
		start := i * 8
		end := start + 8

		var chunk [8]byte
		if end <= needBytes {
			copy(chunk[:], tmp[start:end])
		} else {
			copy(chunk[:], tmp[start:needBytes]) // remaining bytes, rest stay 0
		}
		out[i] = binary.BigEndian.Uint64(chunk[:])
	}

	return out
}

func getPn(IS InternalState, l int) []uint64 {
	if l < 0 || l >= 128 {
		panic("l must satisfy 0 <= l < 128")
	}

	res := make([]uint64, 0, 2)

	if l == 0 {
		return res
	}

	if l <= 64 {
		shift := 64 - l
		mask := ^uint64(0) << shift
		res = append(res, IS.S0&mask)
		return res
	}

	// Case 2: full S0 + part of S1
	res = append(res, IS.S0)

	rem := l - 64
	shift := 64 - rem
	mask := ^uint64(0) << shift
	res = append(res, IS.S1&mask)

	return res
}

func (IS *InternalState) flipBitAtPos(l int) {
	if l < 0 || l >= 320 {
		panic("bit position out of range")
	}

	word := l / 64
	off := l % 64
	mask := uint64(1) << (63 - off)

	switch word {
	case 0:
		IS.S0 ^= mask
	case 1:
		IS.S1 ^= mask
	case 2:
		IS.S2 ^= mask
	case 3:
		IS.S3 ^= mask
	case 4:
		IS.S4 ^= mask
	}
}

func topMask(n int) uint64 {
	// top n bits set (MSB side)
	if n <= 0 {
		return 0
	}
	if n >= 64 {
		return ^uint64(0)
	}
	return ^uint64(0) << (64 - n)
}

// setCnLast overwrites S[0:l-1] with CnLast (C̃_n), MSB-first.
// Assumes l is in [0..128]. CnLast is from your fromByteToUint64(..., l)
// so unused bits (after l) are already zero.
func (IS *InternalState) setCnLast(CnLast []uint64, l int) {
	if l < 0 || l > 128 {
		panic("l must be in [0..128]")
	}
	if l == 0 {
		return
	}

	if l <= 64 {
		// Only touches S0's top l bits.
		if len(CnLast) < 1 {
			panic("CnLast too short")
		}
		m0 := topMask(l)
		IS.S0 = (IS.S0 & ^m0) | (CnLast[0] & m0)
		return
	}

	// l in 65..128: S0 fully overwritten; S1 top (l-64) bits overwritten
	if len(CnLast) < 2 {
		panic("CnLast too short for l>64")
	}

	IS.S0 = CnLast[0]

	rem := l - 64
	m1 := topMask(rem)
	IS.S1 = (IS.S1 & ^m1) | (CnLast[1] & m1)
}

func stringToUint64(s string) ([]uint64, int) {
	b := []byte(s)
	bits := len(b) * 8

	n := (len(b) + 7) / 8
	out := make([]uint64, n)

	for i := 0; i < n; i++ {
		start := i * 8
		end := start + 8
		var chunk [8]byte
		if end <= len(b) {
			copy(chunk[:], b[start:end])
		} else {
			copy(chunk[:], b[start:])
		}
		out[i] = binary.BigEndian.Uint64(chunk[:])

	}
	return out, bits
}

func Uint64ToString(X []uint64, bits int) string {
	if bits < 0 {
		panic("bits must be >= 0")
	}
	if bits > len(X)*64 {
		panic("bits exceeds available words")
	}

	nBytes := (bits + 7) / 8
	if nBytes == 0 {
		return ""
	}

	buf := make([]byte, len(X)*8)

	for i, w := range X {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], w)
	}

	// trim to the meaningful bytes
	buf = buf[:nBytes]

	// if bits isn't a multiple of 8, clear unused low bits in the last byte
	// (keeps the top 'bits%8' bits, MSB-first)
	if rb := bits % 8; rb != 0 {
		mask := byte(0xFF) << uint(8-rb)
		buf[len(buf)-1] &= mask
	}

	return string(buf)
}
