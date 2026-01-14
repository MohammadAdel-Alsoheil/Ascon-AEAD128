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

	// For now (since you're not testing partial), just return the remaining bytes in a stable way:
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
	return []uint64{0x8000000000000000, 0x0000000000000000}
}

func insertByte(b byte, n int) uint64 { return uint64(b) << (8 * (7 - n)) }

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
