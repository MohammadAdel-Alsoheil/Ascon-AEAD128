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
	AD := []uint64{0x0001020304050607, 0x08090A0B0C0D0E0F, 0x0001020304050607, 0x08090A0B0C0D0E0F, 0x0001020304050607, 0x08090A0B0C0D0E0F}
	P := []uint64{0x0001020304050607, 0x08090A0B0C0D0E0F, 0x0001020304050607, 0x08090A0B0C0D0E0F, 0x0001020304050607, 0x08090A0B0C0D0E0F}
	start := time.Now()
	C, T := encrypt(BitString{key, 128}, BitString{nonce, 128}, BitString{AD, 128 * 3}, BitString{P, 128 * 3})
	fmt.Printf("took %v\\n, \n", time.Since(start).Nanoseconds())

	fmt.Printf("C: %x\n", C.Words)
	fmt.Printf("T: %x\n", T.Words)

	start2 := time.Now()
	Pa := decrypt(BitString{key, 128}, BitString{nonce, 128}, BitString{AD, 128 * 3}, C, T)
	fmt.Printf("took %v\\n, \n", time.Since(start2).Nanoseconds())
	fmt.Printf("Pa: %x\n", Pa)
}
