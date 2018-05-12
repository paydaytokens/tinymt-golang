/**
 * Copyright (c) 2011 Mutsuo Saito, Makoto Matsumoto, Hiroshima
 * University and The University of Tokyo. All rights reserved.
 *
 * The 3-clause BSD License is applied to this software, see
 * LICENSE.txt
 *
 */
package tinymt

/*
	https://raw.githubusercontent.com/jj1bdx/tinymtdc-longbatch/master/tinymt64dc/tinymt64dc.255.1048576.txt
	charactristic,                 type, id,    mat1,    mat2,      tmat1,   wght, delta
	8c2c88a7bd41b7051b14736a8d59d79b,64,255,ed561d55,ce58b396,a21ffffffffffffe,63,0
*/

const (
	TINYMT64_MEXP = 127
	TINYMT64_SH0  = 12
	TINYMT64_SH1  = 11
	TINYMT64_SH8  = 8
	TINYMT64_MASK = uint64(0x7fffffffffffffff)
	TINYMT64_MUL  = (1.0 / 9007199254740992.0)
	MIN_LOOP      = 8
)

type TINYMT64_T struct {
	state      [2]uint64
	mat1, mat2 uint32
	tmat       uint64
}

func Tinymt64_init(seed uint64) *TINYMT64_T {
	tmt := &TINYMT64_T{
		mat1: uint32(0xed561d55),
		mat2: uint32(0xce58b396),
		tmat: uint64(0xa21ffffffffffffe),
	}
	tmt.tinymt64_init(seed)
	return tmt
}

func (tmt *TINYMT64_T) tinymt64_init(seed uint64) {
	tmt.state[0] = seed ^ (uint64(tmt.mat1) << 32)
	tmt.state[1] = uint64(tmt.mat2) ^ tmt.tmat
	for i := 1; i < MIN_LOOP; i++ {
		tmt.state[i&1] ^= uint64(i) + uint64(6364136223846793005)*(tmt.state[(i-1)&1]^(tmt.state[(i-1)&1]>>62))
	}
	tmt.period_certification()
}

func (tmt *TINYMT64_T) period_certification() {
	if (tmt.state[0]&TINYMT64_MASK) == 0 && tmt.state[1] == 0 {
		tmt.state[0] = 'T'
		tmt.state[1] = 'M'
	}
}

func (tmt *TINYMT64_T) tinymt64_next_state() {
	var x uint64

	tmt.state[0] &= TINYMT64_MASK
	x = tmt.state[0] ^ tmt.state[1]
	x ^= x << TINYMT64_SH0
	x ^= x >> 32
	x ^= x << 32
	x ^= x << TINYMT64_SH1
	tmt.state[0] = tmt.state[1]
	tmt.state[1] = x
	tmt.state[0] ^= uint64(-(int64(x & uint64(1))) & int64(tmt.mat1))
	tmt.state[1] ^= uint64(-(int64(x & 1))) & (uint64(tmt.mat2) << 32)
}

func (tmt *TINYMT64_T) tinymt64_temper() uint64 {
	var x uint64
	// linearity check
	// x = tmt.state[0] ^ tmt.state[1]
	x = tmt.state[0] + tmt.state[1]
	x ^= tmt.state[0] >> TINYMT64_SH8
	x ^= uint64(-(int64(x & 1))) & tmt.tmat
	return x
}

func (tmt *TINYMT64_T) Generate() uint64 {
	tmt.tinymt64_next_state()
	return tmt.tinymt64_temper()
}
