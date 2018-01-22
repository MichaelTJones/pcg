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
func TestSanity64(t *testing.T) {
	result := NewPCG64().Seed(1, 1, 1, 2).Random()
	expect := uint64(1107300197865787281)
	if result != expect {
		t.Errorf("NewPCG64().Seed(1, 1, 1, 2).Random() is %q; want %q", result, expect)
	}
}

var sumTests64 = []struct {
	state    uint64 // PCG seed value for state
	sequence uint64 // PCG seed value for sequence
	count    int    // number of values to sum
	sum      uint64 // sum of the first count values
}{
	{1, 1, 10, 8034187309725975364},
	{1, 1, 100, 14328956917741108809},
	{1, 1, 1000, 15814724732829753998},
	{1, 1, 10000, 8547922387302793844},
}

// Are the sums of the first few values consistent with expectation?
func TestSum64(t *testing.T) {
	for i, a := range sumTests64 {
		pcg := NewPCG64().Seed(a.state, a.state, a.sequence, a.sequence+1)
		sum := uint64(0)
		for j := 0; j < a.count; j++ {
			sum += uint64(pcg.Random())
		}
		if sum != a.sum {
			t.Errorf("#%d, sum of first %d values = %d; want %d", i, a.count, sum, a.sum)
		}
	}
}

const count64 = 256

// Does advancing work?
func TestAdvance64(t *testing.T) {
	pcg := NewPCG64().Seed(1, 1, 1, 2)
	values := make([]uint64, count64)
	for i := range values {
		values[i] = pcg.Random()
	}

	for skip := 1; skip < count64; skip++ {
		pcg.Seed(1, 1, 1, 2)
		pcg.Advance(uint64(skip))
		result := pcg.Random()
		expect := values[skip]
		if result != expect {
			t.Errorf("Advance(%d) is %d; want %d", skip, result, expect)
		}
	}
}

// Does retreating work?
func TestRetreat64(t *testing.T) {
	pcg := NewPCG64().Seed(1, 1, 1, 2)
	expect := pcg.Random()

	for skip := 1; skip < count64; skip++ {
		pcg.Seed(1, 1, 1, 2)
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

// Measure the time it takes to generate a 64-bit generator
func BenchmarkNew64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPCG64().Seed(1, 1, 1, 2)
	}
}

// Measure the time it takes to generate random values
func BenchmarkRandom64(b *testing.B) {
	b.StopTimer()
	pcg := NewPCG64().Seed(1, 1, 1, 2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = pcg.Random()
	}
}

// Measure the time it takes to generate bounded random values
func BenchmarkBounded64(b *testing.B) {
	b.StopTimer()
	pcg := NewPCG64().Seed(1, 1, 1, 2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = pcg.Bounded(uint64(i) & 0xff) // 0..255
		// _ = pcg.Bounded(6)             // roll of die
		// _ = pcg.Bounded(52)            // deck of cards
		// _ = pcg.Bounded(365)           // day of year
	}
}

//
// EXAMPLES
//

func ExampleReport64() {
	// Print report
	rng := NewPCG64().Seed(42, 42, 54, 54)

	fmt.Printf("pcg32x2 random:\n"+
		"      -  result:      64-bit unsigned int (uint64)\n"+
		"      -  period:      2^64   (* 2^63 streams)\n"+
		"      -  state type:  PGC64 (%d bytes)\n"+
		"      -  output func: XSH-RR\n"+
		"\n",
		unsafe.Sizeof(*rng))

	for round := 1; round <= 5; round++ {
		fmt.Printf("Round %d:\n", round)

		/* Make some 64-bit numbers */
		fmt.Printf("  64bit:")
		for i := 0; i < 6; i++ {
			if i > 0 && i%3 == 0 {
				fmt.Printf("\n\t")
			}
			fmt.Printf(" 0x%016x", rng.Random())
		}
		fmt.Println()

		fmt.Printf("  Again:")
		// rng.Advance(-6 & (1<<64 - 1))
		// rng.Advance(1<<64 - 6)
		rng.Retreat(6)
		for i := 0; i < 6; i++ {
			if i > 0 && i%3 == 0 {
				fmt.Printf("\n\t")
			}
			fmt.Printf(" 0x%016x", rng.Random())
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
		for i := uint64(CARDS); i > 1; i-- {
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
	// pcg32x2 random:
	//       -  result:      64-bit unsigned int (uint64)
	//       -  period:      2^64   (* 2^63 streams)
	//       -  state type:  PGC64 (16 bytes)
	//       -  output func: XSH-RR
	//
	// Round 1:
	//   64bit: 0x1a410f65a15c02b7 0xe0b09a537b47f409 0x11fba8acba1d3330
	// 	 0x452993e983d2f293 0x36082c12bfa4784b 0xf5934191cbed606e
	//   Again: 0x1a410f65a15c02b7 0xe0b09a537b47f409 0x11fba8acba1d3330
	// 	 0x452993e983d2f293 0x36082c12bfa4784b 0xf5934191cbed606e
	//   Coins: HHTTTHTHHHTHTTTHHHHHTTTHHHTHTHTHTTHTTTHHHHHHTTTTHHTTTTTHTTTTTTTHT
	//   Rolls: 3 2 1 3 2 2 3 6 4 5 6 4 5 5 5 2 1 1 1 3 3 3 6 1 3 6 6 6 2 2 6 3 5
	//   Cards: 5c 9s Th 5s 2h 7d 5d Tc Qh 8h 9d 4s 3d 8d 7c 7s Js Qc 3h Jc 2d Jd
	// 	 Ad 6d 3c Kh 7h Qd 9h Qs 9c 6s As 4c Td 6c Kd Ks 8s 8c 2s 4d 6h Ah
	// 	 3s 2c 4h 5h Ts Ac Jh Kc
	// Round 2:
	//   64bit: 0x2c59d33474ab93ad 0x63e4d9a81c1da000 0xd80ced60494ff896
	// 	 0xabb62df034462f2f 0x6ceea41fd308a3e5 0x7c561bd70fa83bab
	//   Again: 0x2c59d33474ab93ad 0x63e4d9a81c1da000 0xd80ced60494ff896
	// 	 0xabb62df034462f2f 0x6ceea41fd308a3e5 0x7c561bd70fa83bab
	//   Coins: HHHHHHHHHHTHHHTHTHTHTHTTTTHHTTTHHTHHTHTTHHTTTHHHHHHTHTTHTHTTTTTTT
	//   Rolls: 5 5 3 5 5 6 4 5 1 4 6 4 2 1 6 2 6 6 3 6 5 4 5 6 3 4 4 6 3 3 5 2 2
	//   Cards: 6d 9d 3c 7h Qc Ts Js 3h Th 4s 5d 5c 2d 3d 4c 6c Ad Kd 9c 6h Qs 2c
	// 	 Td 8h 5s 5h 8s 7c 7d Jh Kh Qd 4d Jc Kc 3s Ah 7s Jd 9h 4h 8d Tc 6s
	// 	 2h Qh Ks 9s Ac 2s As 8c
	// Round 3:
	//   64bit: 0x7088079939af5f9f 0x31198a8804196b18 0x26859955c3c3eb28
	// 	 0x822c3b19c076c60c 0x1c41c1e0c693e135 0x7255b24df8f63932
	//   Again: 0x7088079939af5f9f 0x31198a8804196b18 0x26859955c3c3eb28
	// 	 0x822c3b19c076c60c 0x1c41c1e0c693e135 0x7255b24df8f63932
	//   Coins: HTTHHTTTTTHTTHHHTHTTHHTTHTHHTHTHTTTTHHTTTHHTHHTTHTTHHHTHHHTHTTTHT
	//   Rolls: 3 3 5 3 4 4 2 3 1 1 5 3 4 6 5 2 1 2 2 4 1 3 3 2 6 2 4 2 2 2 5 5 4
	//   Cards: Ac 5h Ts Kd 2d 6d 5d 6h 3d 8c Jd 5s 9h 4h Tc 9d Qh 7d Ad 4s 5c 2c
	// 	 Ah Qc Jh Th 6c 3s 8h 7s 7c Js 8s 4d Kh 9s Qs 7h Ks Qd 2h Jc 9c Kc
	// 	 As 4c 8d Td 3h 6s 3c 2s
	// Round 4:
	//   64bit: 0xe2d9967955ce6851 0x339ab6aa97a7726d 0x7638b42017e10815
	// 	 0x742b519858007d43 0x6083910e962fb148 0xdcb7a611b9bb55bd
	//   Again: 0xe2d9967955ce6851 0x339ab6aa97a7726d 0x7638b42017e10815
	// 	 0x742b519858007d43 0x6083910e962fb148 0xdcb7a611b9bb55bd
	//   Coins: HHTHHTTTTHTHHHHHTTHHHTTTHHTHTHTHTHHTTHTHHHHHHTHHTHHTHHTTTTHHTHHTT
	//   Rolls: 6 6 3 6 1 4 6 2 2 2 4 1 2 3 3 5 5 2 2 2 6 4 2 3 2 6 4 5 6 3 1 2 5
	//   Cards: 7c 4c Qh 2h 7s 3s 4h 2c Ad 8c 5s 9h 2d 4s Td 6d 2s 8h Ks 6h 9d 3h
	// 	 Ac Th 7d 5c Jd 8s 5d 6c Qs 3c 5h Tc As Qc Kd Ts 9s 3d Qd Ah 4d 8d
	// 	 Kh Js Jc Jh Kc 6s 7h 9c
	// Round 5:
	//   64bit: 0x7f85ce72fcef7cd6 0xddfc86301b488b5a 0xc38e77dbd0daf7ea
	// 	 0x27df76751d9a70f7 0x77011f4d241a37cf 0xbe702c2b9a3857b7
	//   Again: 0x7f85ce72fcef7cd6 0xddfc86301b488b5a 0xc38e77dbd0daf7ea
	// 	 0x27df76751d9a70f7 0x77011f4d241a37cf 0xbe702c2b9a3857b7
	//   Coins: HHHHTHHTTHTTHHHTTTHHTHTHTTTTHTTHTHTTTHHHTHTHTTHTTHTHHTHTHHHTHTHTT
	//   Rolls: 3 4 3 6 2 1 5 1 5 6 3 4 2 5 2 4 3 6 5 3 2 5 2 6 6 2 1 4 6 2 5 3 2
	//   Cards: 6c 2d 9s Ah 2s As 4h 6h Kd 4d 3d Js 8d Qs Jh Qd 6s 9d 3h Qh Th 5c
	// 	 2h Kc 7c 5s 7s 5d Ks Ad 7h Tc Jc 7d 3s 8c 8s 3c 6d Ac 9c 2c Td Qc
	// 	 5h 9h 8h Kh Ts Jd 4c 4s
}
