package main

import "math/bits"

func (IS *InternalState) permute(rnd int) {

	for i := 0; i < rnd; i++ {
		IS.Pc(i, rnd)
		IS.Ps()
		IS.Pl()
	}
}

func (IS *InternalState) Pc(i int, rnd int) {
	var constants = []uint64{
		0x000000000000003c,
		0x000000000000002d,
		0x000000000000001e,
		0x000000000000000f,
		0x00000000000000f0,
		0x00000000000000e1,
		0x00000000000000d2,
		0x00000000000000c3,
		0x00000000000000b4,
		0x00000000000000a5,
		0x0000000000000096,
		0x0000000000000087,
		0x0000000000000078,
		0x0000000000000069,
		0x000000000000005a,
		0x000000000000004b,
	}

	IS.S2 ^= constants[16-rnd+i]
}

func (IS *InternalState) Ps() {
	x0 := IS.S0
	x1 := IS.S1
	x2 := IS.S2
	x3 := IS.S3
	x4 := IS.S4

	x0 ^= x4
	x4 ^= x3
	x2 ^= x1

	t0 := (^x0) & x1
	t1 := (^x1) & x2
	t2 := (^x2) & x3
	t3 := (^x3) & x4
	t4 := (^x4) & x0

	x0 ^= t1
	x1 ^= t2
	x2 ^= t3
	x3 ^= t4
	x4 ^= t0

	x1 ^= x0
	x0 ^= x4
	x3 ^= x2
	x2 = ^x2

	IS.S0 = x0
	IS.S1 = x1
	IS.S2 = x2
	IS.S3 = x3
	IS.S4 = x4
}

func (IS *InternalState) Pl() {
	x01 := bits.RotateLeft64(IS.S0, 45)
	x02 := bits.RotateLeft64(IS.S0, 36)

	x11 := bits.RotateLeft64(IS.S1, 3)
	x12 := bits.RotateLeft64(IS.S1, 25)

	x21 := bits.RotateLeft64(IS.S2, 63)
	x22 := bits.RotateLeft64(IS.S2, 58)

	x31 := bits.RotateLeft64(IS.S3, 54)
	x32 := bits.RotateLeft64(IS.S3, 47)

	x41 := bits.RotateLeft64(IS.S4, 57)
	x42 := bits.RotateLeft64(IS.S4, 23)

	IS.S0 = IS.S0 ^ x01 ^ x02
	IS.S1 = IS.S1 ^ x11 ^ x12
	IS.S2 = IS.S2 ^ x21 ^ x22
	IS.S3 = IS.S3 ^ x31 ^ x32
	IS.S4 = IS.S4 ^ x41 ^ x42
}
