package pcg

// Copyright 2018 Michael T. Jones
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance 
// with the License. You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed 
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for 
// the specific language governing permissions and limitations under the License.

import (
	"fmt"
	"testing"
	"unsafe"
)

//
// TESTS
//

// Basic sanity test: is first known value determined properly?
func TestSanity32(t *testing.T) {
	result := NewPCG32().Seed(1, 1).Random()
	expect := uint32(3380776849)
	if result != expect {
		t.Errorf("NewPCG32().Seed(1, 1).Random() is %d; want %d", result, expect)
	}
}

var sumTests32 = []struct {
	state    uint64 // PCG seed value for state
	sequence uint64 // PCG seed value for sequence
	count    int    // number of values to sum
	sum      uint64 // sum of the first count values
}{
	{1, 1, 10, 23876797252},
	{1, 1, 100, 214946922057},
	{1, 1, 1000, 2103475364494},
	{1, 1, 10000, 21493472477812},
}

// Are the sums of the first few values consistent with expectation?
func TestSum32(t *testing.T) {
	for i, a := range sumTests32 {
		pcg := NewPCG32().Seed(a.state, a.sequence)
		sum := uint64(0)
		for j := 0; j < a.count; j++ {
			sum += uint64(pcg.Random())
		}
		if sum != a.sum {
			t.Errorf("#%d, sum of first %d values = %d; want %d", i, a.count, sum, a.sum)
		}
	}
}

const count32 = 256

// Does advancing work?
func TestAdvance32(t *testing.T) {
	pcg := NewPCG32().Seed(1, 1)
	values := make([]uint32, count32)
	for i := range values {
		values[i] = pcg.Random()
	}

	for skip := 1; skip < count32; skip++ {
		pcg.Seed(1, 1)
		pcg.Advance(uint64(skip))
		result := pcg.Random()
		expect := values[skip]
		if result != expect {
			t.Errorf("Advance(%d) is %d; want %d", skip, result, expect)
		}
	}
}

// Does retreating work?
func TestRetreat32(t *testing.T) {
	pcg := NewPCG32().Seed(1, 1)
	expect := pcg.Random()

	for skip := 1; skip < count32; skip++ {
		pcg.Seed(1, 1)
		for i := 0; i < skip; i++ {
			_ = pcg.Random()
		}
		pcg.Retreat(uint64(skip))
		result := pcg.Random()
		if result != expect {
			t.Errorf("Retreat(%d) is %d; want %d", skip, result, expect)
		}
	}
}

//
// BENCHMARKS
//

// Measure the time it takes to generate a 32-bit generator
func BenchmarkNew32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPCG32().Seed(1, 1)
	}
}

// Measure the time it takes to generate random values
func BenchmarkRandom32(b *testing.B) {
	b.StopTimer()
	pcg := NewPCG32().Seed(1, 1)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = pcg.Random()
	}
}

// Measure the time it takes to generate bounded random values
func BenchmarkBounded32(b *testing.B) {
	b.StopTimer()
	pcg := NewPCG32().Seed(1, 1)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = pcg.Bounded(uint32(i) & 0xff) // 0..255
		// _ = pcg.Bounded(6)             // roll of die
		// _ = pcg.Bounded(52)            // deck of cards
		// _ = pcg.Bounded(365)           // day of year
	}
}

//
// EXAMPLES
//

func ExampleReport32() {
	// Print report
	rng := NewPCG32().Seed(42, 54)

	fmt.Printf("pcg32 random:\n"+
		"      -  result:      32-bit unsigned int (uint32)\n"+
		"      -  period:      2^64   (* 2^63 streams)\n"+
		"      -  state type:  PGC32 (%d bytes)\n"+
		"      -  output func: XSH-RR\n"+
		"\n",
		unsafe.Sizeof(*rng))

	for round := 1; round <= 5; round++ {
		fmt.Printf("Round %d:\n", round)

		/* Make some 32-bit numbers */
		fmt.Printf("  32bit:")
		for i := 0; i < 6; i++ {
			fmt.Printf(" 0x%08x", rng.Random())
		}
		fmt.Println()

		fmt.Printf("  Again:")
		// rng.Advance(-6 & (1<<64 - 1))
		// rng.Advance(1<<64 - 6)
		rng.Retreat(6)
		for i := 0; i < 6; i++ {
			fmt.Printf(" 0x%08x", rng.Random())
		}
		fmt.Println()

		/* Toss some coins */
		fmt.Printf("  Coins: ")
		for i := 0; i < 65; i++ {
			fmt.Printf("%c", "TH"[rng.Bounded(2)])
		}
		fmt.Println()

		/* Roll some dice */
		fmt.Printf("  Rolls:")
		for i := 0; i < 33; i++ {
			fmt.Printf(" %d", rng.Bounded(6)+1)
		}
		fmt.Println()

		/* Deal some cards */
		const (
			SUITS = 4
			CARDS = 52
		)
		var cards [CARDS]int
		for i := range cards {
			cards[i] = i
		}
		for i := uint32(CARDS); i > 1; i-- {
			chosen := rng.Bounded(i)
			cards[chosen], cards[i-1] = cards[i-1], cards[chosen]
		}

		fmt.Printf("  Cards:")
		for i, c := range cards {
			fmt.Printf(" %c%c", "A23456789TJQK"[c/SUITS], "hcds"[c%SUITS])
			if (i+1)%22 == 0 {
				fmt.Printf("\n\t")
			}
		}
		fmt.Println()

	}
	// Output:
	// pcg32 random:
	//       -  result:      32-bit unsigned int (uint32)
	//       -  period:      2^64   (* 2^63 streams)
	//       -  state type:  PGC32 (16 bytes)
	//       -  output func: XSH-RR
	//
	// Round 1:
	//   32bit: 0xa15c02b7 0x7b47f409 0xba1d3330 0x83d2f293 0xbfa4784b 0xcbed606e
	//   Again: 0xa15c02b7 0x7b47f409 0xba1d3330 0x83d2f293 0xbfa4784b 0xcbed606e
	//   Coins: HHTTTHTHHHTHTTTHHHHHTTTHHHTHTHTHTTHTTTHHHHHHTTTTHHTTTTTHTTTTTTTHT
	//   Rolls: 3 4 1 1 2 2 3 2 4 3 2 4 3 3 5 2 3 1 3 1 5 1 4 1 5 6 4 6 6 2 6 3 3
	//   Cards: Qd Ks 6d 3s 3d 4c 3h Td Kc 5c Jh Kd Jd As 4s 4h Ad Th Ac Jc 7s Qs
	// 	 2s 7h Kh 2d 6c Ah 4d Qh 9h 6s 5s 2c 9c Ts 8d 9s 3c 8c Js 5d 2h 6h
	// 	 7d 8s 9d 5h 8h Qc 7c Tc
	// Round 2:
	//   32bit: 0x74ab93ad 0x1c1da000 0x494ff896 0x34462f2f 0xd308a3e5 0x0fa83bab
	//   Again: 0x74ab93ad 0x1c1da000 0x494ff896 0x34462f2f 0xd308a3e5 0x0fa83bab
	//   Coins: HHHHHHHHHHTHHHTHTHTHTHTTTTHHTTTHHTHHTHTTHHTTTHHHHHHTHTTHTHTTTTTTT
	//   Rolls: 5 1 1 3 3 2 4 5 3 2 2 6 4 3 2 4 2 4 3 2 3 6 3 2 3 4 2 4 1 1 5 4 4
	//   Cards: 7d 2s 7h Td 8s 3c 3d Js 2d Tc 4h Qs 5c 9c Th 2c Jc Qd 9d Qc 7s 3s
	// 	 5s 6h 4d Jh 4c Ac 4s 5h 5d Kc 8h 8d Jd 9s Ad 6s 6c Kd 2h 3h Kh Ts
	// 	 Qh 9h 6d As 7c Ks Ah 8c
	// Round 3:
	//   32bit: 0x39af5f9f 0x04196b18 0xc3c3eb28 0xc076c60c 0xc693e135 0xf8f63932
	//   Again: 0x39af5f9f 0x04196b18 0xc3c3eb28 0xc076c60c 0xc693e135 0xf8f63932
	//   Coins: HTTHHTTTTTHTTHHHTHTTHHTTHTHHTHTHTTTTHHTTTHHTHHTTHTTHHHTHHHTHTTTHT
	//   Rolls: 5 1 5 3 2 2 4 5 3 3 1 3 4 6 3 2 3 4 2 2 3 1 5 2 4 6 6 4 2 4 3 3 6
	//   Cards: Kd Jh Kc Qh 4d Qc 4h 9d 3c Kh Qs 8h 5c Jd 7d 8d 3h 7c 8s 3s 2h Ks
	// 	 9c 9h 2c 8c Ad 7s 4s 2s 5h 6s 4c Ah 7h 5s Ac 3d 5d Qd As Tc 6h 9s
	// 	 2d 6c 6d Td Jc Ts Th Js
	// Round 4:
	//   32bit: 0x55ce6851 0x97a7726d 0x17e10815 0x58007d43 0x962fb148 0xb9bb55bd
	//   Again: 0x55ce6851 0x97a7726d 0x17e10815 0x58007d43 0x962fb148 0xb9bb55bd
	//   Coins: HHTHHTTTTHTHHHHHTTHHHTTTHHTHTHTHTHHTTHTHHHHHHTHHTHHTHHTTTTHHTHHTT
	//   Rolls: 6 6 3 2 3 4 2 6 4 2 6 3 2 3 5 5 3 4 4 6 6 2 6 5 4 4 6 1 6 1 3 6 5
	//   Cards: Qd 8h 5d 8s 8d Ts 7h Th Qs Js 7s Kc 6h 5s 4d Ac Jd 7d 7c Td 2c 6s
	// 	 5h 6d 3s Kd 9s Jh Kh As Ah 9h 3c Qh 9c 2d Tc 9d 2s 3d Ks 4h Qc Ad
	// 	 Jc 8c 2h 3h 4s 4c 5c 6c
	// Round 5:
	//   32bit: 0xfcef7cd6 0x1b488b5a 0xd0daf7ea 0x1d9a70f7 0x241a37cf 0x9a3857b7
	//   Again: 0xfcef7cd6 0x1b488b5a 0xd0daf7ea 0x1d9a70f7 0x241a37cf 0x9a3857b7
	//   Coins: HHHHTHHTTHTTHHHTTTHHTHTHTTTTHTTHTHTTTHHHTHTHTTHTTHTHHTHTHHHTHTHTT
	//   Rolls: 5 4 1 2 6 1 3 1 5 6 3 6 2 1 4 4 5 2 1 5 6 5 6 4 4 4 5 2 6 4 3 5 6
	//   Cards: 4d 9s Qc 9h As Qs 7s 4c Kd 6h 6s 2c 8c 5d 7h 5h Jc 3s 7c Jh Js Ks
	// 	 Tc Jd Kc Th 3h Ts Qh Ad Td 3c Ah 2d 3d 5c Ac 8s 5s 9c 2h 6c 6d Kh
	// 	 Qd 8d 7d 2s 8h 4h 9d 4s
}
