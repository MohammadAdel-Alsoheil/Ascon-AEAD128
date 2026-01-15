package main

import (
	"fmt"
	"time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	// reference https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-232.pdf

	key := []uint64{0x0001020304050607, 0x08090A0B0C0D0E0F}
	nonce := []uint64{0x0001020304050607, 0x08090A0B0C0D0E0F}
	AD, ADBits := stringToUint64("HelloHelloHelloHHBBBB")   // 18
	P, Pbits := stringToUint64("Hello, I love you so much") //24

	start := time.Now()
	C, T := encrypt(BitString{key, 128}, BitString{nonce, 128}, BitString{AD, ADBits}, BitString{P, Pbits})
	fmt.Printf("took %v\\n, \n", time.Since(start).Nanoseconds())

	fmt.Printf("CIPHER is : %s\n", Uint64ToString(C.Words, C.Bits))
	fmt.Printf("TAG is: %s\n", Uint64ToString(T.Words, T.Bits))

	start2 := time.Now()
	Pa := decrypt(BitString{key, 128}, BitString{nonce, 128}, BitString{AD, ADBits}, C, T)
	fmt.Printf("took %v\\n, \n", time.Since(start2).Nanoseconds())
	fmt.Printf("Deciphered PlainText is : %s \n", Uint64ToString(Pa.Words, Pa.Bits))
}
