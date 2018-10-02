package main

//IsaacCipher is the struct containing our state
type IsaacCipher struct {
	randResult [256]int
	valuesLeft int
	mm [256]int
	aa int
	bb int
	cc int
}

func (i *IsaacCipher)generateMoreResults() {
	i.cc++
	i.bb += i.cc
	for z := 0; z < 256; z++ {
		fb := i.mm[z]
		switch (z % 4) { //usually z%4 was z&3
		case 0:
			i.aa = i.aa ^ (i.aa << 13)
		case 1:
			i.aa = i.aa ^ (i.aa >> 6)
		case 2:
			i.aa = i.aa ^ (i.aa << 2)
		case 3:
			i.aa = i.aa ^ (i.aa >> 16)
		}
		i.aa = i.mm[z ^ 128] + i.aa
		i.mm[z] = i.mm[(z >> 2) & 0xff] + i.aa + i.bb
		u := i.mm[z]
		i.bb = i.mm[(u >> 10) & 0xff] + fb
		i.randResult[z] = i.bb
	}
	i.valuesLeft = 256
}

func (i *IsaacCipher)Init(seed []int) {
	i.cc = 0
	i.bb = i.cc
	i.aa = i.bb

	initState := make([]int, 8)
	for o := 0; o < len(initState); o++ {
		initState[o] = 0x9e3779b9 //golden ratio
	}
	/*initState[0] = 0x9
	initState[1] = 0xe
	initState[2] = 0x3 
	initState[3] = 0x7
	initState[4] = 0x7
	initState[5] = 0x9
	initState[6] = 0xb
	initState[7] = 0x9*/
	for z := 0; z < 4; z++ {
		i.mix(initState)
	}
	for y := 0; y < 256; y += 8 {
		if seed != nil && len(seed) == 256 {
			for w := 0; w < 8; w++ { //seed must be 256 otherwise this will crash
				initState[w] += seed[y+w]
			}
		}
		i.mix(initState)
		for w := 0; w < 8; w++ {
			i.mm[y+w] = initState[w]
		}
	}
	if seed != nil {
		for q := 0; q < 256; q += 8 {
			for w := 0; w < 8; w++ {
				initState[w] += i.mm[q+w]
			}
			i.mix(initState)
			for w := 0; w < 8; w++ {
				i.mm[q+w] = initState[w]
			}
		}
	}
	i.valuesLeft = 0
}

func (i *IsaacCipher)mix(s []int) {
	s[0] ^= s[1] << 11
	s[3] += s[0]
	s[1] += s[2]
	s[1] ^= s[2] >> 2
	s[4] += s[1]
	s[2] += s[3]
	s[2] ^= s[3] << 8
	s[5] += s[2]
	s[3] += s[4]
	s[3] ^= s[4] >> 16
	s[6] += s[3]
	s[4] += s[5]
	s[4] ^= s[5] << 10
	s[7] += s[4]
	s[5] += s[6]
	s[5] ^= s[6] >> 4
	s[0] += s[5]
	s[6] += s[7]
	s[6] ^= s[7] << 8
	s[1] += s[6]
	s[7] += s[0]
	s[7] ^= s[0] >> 9
	s[2] += s[7]
	s[0] += s[1]
}

func (i *IsaacCipher)NextInt() int {
	if i.valuesLeft == 0 {
		i.generateMoreResults()
	}
	i.valuesLeft--
	return i.randResult[i.valuesLeft]
}
