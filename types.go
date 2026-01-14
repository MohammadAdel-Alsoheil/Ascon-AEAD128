package main

type InternalState struct {
	S0 uint64
	S1 uint64
	S2 uint64
	S3 uint64
	S4 uint64
}

type BitString struct {
	Words []uint64
	Bits  int
}
