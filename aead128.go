package main

var IV uint64 = 0x00001000808c0001

func encrypt(key BitString, Nonce BitString, AD BitString, P BitString) (BitString, BitString) {

	IS := InternalState{
		S0: IV,
		S1: key.Words[0],
		S2: key.Words[1],
		S3: Nonce.Words[0],
		S4: Nonce.Words[1],
	}
	IS.permute(12)
	IS.S0 = IS.S0 ^ 0
	IS.S1 = IS.S1 ^ 0
	IS.S2 = IS.S2 ^ 0
	IS.S3 = IS.S3 ^ key.Words[0]
	IS.S4 = IS.S4 ^ key.Words[1]

	if AD.Bits > 0 {
		ADs, Am := parse(AD, 128)
		ADs = append(ADs, pad(Am, 128))
		for i := 0; i < len(ADs); i++ {
			IS.S0 = IS.S0 ^ ADs[i][0]
			IS.S1 = IS.S1 ^ ADs[i][1]
			IS.permute(8)
		}
	}
	IS.S0 = IS.S0 ^ 0
	IS.S1 = IS.S1 ^ 0
	IS.S2 = IS.S2 ^ 0
	IS.S3 = IS.S3 ^ 0
	IS.S4 = IS.S4 ^ 1

	Ps, Pn := parse(P, 128)
	l := P.Bits % 128
	C := BitString{
		Words: make([]uint64, 0),
		Bits:  0,
	}
	for i := 0; i < len(Ps); i++ {
		IS.S0 = IS.S0 ^ Ps[i][0]
		IS.S1 = IS.S1 ^ Ps[i][1]
		C.Words = append(C.Words, IS.S0, IS.S1)
		C.Bits += 128
		IS.permute(8)
	}
	temp := pad(Pn, 128)
	IS.S0 = IS.S0 ^ temp[0]
	IS.S1 = IS.S1 ^ temp[1]

	C.Words = append(C.Words, getPn(IS, l)...)
	C.Bits += l

	IS.S0 = IS.S0 ^ 0
	IS.S1 = IS.S1 ^ 0
	IS.S2 = IS.S2 ^ key.Words[0]
	IS.S3 = IS.S3 ^ key.Words[1]
	IS.S4 = IS.S4 ^ 0

	IS.permute(12)
	T := BitString{
		make([]uint64, 0), 128,
	}

	T.Words = append(T.Words, IS.S3)
	T.Words = append(T.Words, IS.S4)
	T.Words[0] = T.Words[0] ^ key.Words[0]
	T.Words[1] = T.Words[1] ^ key.Words[1]

	return C, T
}

func decrypt(key BitString, Nonce BitString, AD BitString, C BitString, T BitString) BitString {

	IS := InternalState{
		S0: IV,
		S1: key.Words[0],
		S2: key.Words[1],
		S3: Nonce.Words[0],
		S4: Nonce.Words[1],
	}
	IS.permute(12)
	IS.S0 = IS.S0 ^ 0
	IS.S1 = IS.S1 ^ 0
	IS.S2 = IS.S2 ^ 0
	IS.S3 = IS.S3 ^ key.Words[0]
	IS.S4 = IS.S4 ^ key.Words[1]

	if AD.Bits > 0 {
		ADs, Am := parse(AD, 128)
		ADs = append(ADs, pad(Am, 128))
		for i := 0; i < len(ADs); i++ {
			IS.S0 = IS.S0 ^ ADs[i][0]
			IS.S1 = IS.S1 ^ ADs[i][1]
			IS.permute(8)
		}
	}
	IS.S0 = IS.S0 ^ 0
	IS.S1 = IS.S1 ^ 0
	IS.S2 = IS.S2 ^ 0
	IS.S3 = IS.S3 ^ 0
	IS.S4 = IS.S4 ^ 1

	Cs, Cn := parse(C, 128)
	l := C.Bits % 128
	P := BitString{
		Words: make([]uint64, 0),
		Bits:  0,
	}
	for i := 0; i < len(Cs); i++ {
		IS.S0 = IS.S0 ^ Cs[i][0]
		IS.S1 = IS.S1 ^ Cs[i][1]
		P.Words = append(P.Words, IS.S0, IS.S1)
		P.Bits += 128
		IS.S0 = Cs[i][0]
		IS.S1 = Cs[i][1]
		IS.permute(8)
	}
	temp := getPn(IS, l) // S[0:l-1]
	CnLast := fromByteToUint64(Cn, l)
	PnLast := make([]uint64, 0)
	for i := 0; i < len(temp); i++ {
		PnLast = append(PnLast, CnLast[i]^temp[i])
	}
	IS.flipBitAtPos(l)
	IS.setCnLast(CnLast, l)

	IS.S0 = IS.S0 ^ 0
	IS.S1 = IS.S1 ^ 0
	IS.S2 = IS.S2 ^ key.Words[0]
	IS.S3 = IS.S3 ^ key.Words[1]
	IS.S4 = IS.S4 ^ 0
	IS.permute(12)

	Tprime := make([]uint64, 0)

	Tprime = append(Tprime, IS.S3^key.Words[0], IS.S4^key.Words[1])

	if Tprime[0] == T.Words[0] && Tprime[1] == T.Words[1] {
		P.Words = append(P.Words, PnLast...)
		P.Bits += l
		return P
	}
	panic("Integrity check failed")

}
